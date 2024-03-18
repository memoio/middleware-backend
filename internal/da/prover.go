package da

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/kzg"
	"github.com/memoio/backend/api"
	proof "github.com/memoio/go-did/file-proof"
)

type DataAccessProver struct {
	proofInstance      proof.ProofInstance
	selectedFileNumber int64
	provingKey         kzg.ProvingKey
	interval           time.Duration
	period             time.Duration
	respondTime        int64
	last               int64
}

func NewDataAccessProver(chain string, sk *ecdsa.PrivateKey) (*DataAccessProver, error) {
	instance, err := proof.NewProofInstance(sk, chain)
	if err != nil {
		return nil, err
	}

	info, err := instance.GetSettingInfo()
	if err != nil {
		return nil, err
	}

	return &DataAccessProver{
		proofInstance:      *instance,
		selectedFileNumber: int64(info.ChalSum),
		provingKey:         DefaultSRS.Pk,
		interval:           time.Duration(info.Interval) * time.Second,
		period:             time.Duration(info.Period) * time.Second,
		respondTime:        int64(info.RespondTime),
		last:               0,
	}, nil
}

func (p *DataAccessProver) ProveDataAccess() {
	var lastRnd fr.Element
	var nowRnd fr.Element
	var proveSuccess bool
	for p.last == 0 {
		time.Sleep(5 * time.Second)
		_, _, lastTime, err := p.proofInstance.GetVerifyInfo()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		p.last = lastTime.Int64()
	}

	for {
		wait, nextTime := p.calculateWatingTime()
		time.Sleep(wait)

		p.resetChallengeStatus()

		var err error
		var lock bool
		start := time.Now()
		logger.Info("start prove")
		lastRnd = nowRnd
		nowRnd, lock, err = p.generateRND()
		if err != nil {
			logger.Error(err.Error())
			proveSuccess = false
			continue
		}

		if lastRnd.Cmp(&nowRnd) == 0 && proveSuccess {
			logger.Error("rnd shouldn't be the same as before")
			continue
		}
		if !lock {
			logger.Error("verify should be locked after generate rnd")
			continue
		}

		commits, proofs, err := p.selectFiles(nowRnd)
		if err != nil {
			logger.Error(err.Error())
			proveSuccess = false
			continue
		}

		err = p.proveToContract(commits, proofs, nowRnd)
		if err != nil {
			logger.Error(err.Error())
			proveSuccess = false
			continue
		}
		logger.Infof("end prove, using: %fs", time.Since(start).Seconds())

		proveSuccess = true

		start = time.Now()
		logger.Info("start response chanllenge")
		p.last = nextTime

		err = p.responseChallenge(commits)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		logger.Infof("end response challenge, using: %fs", time.Since(start).Seconds())
	}
}

func (p *DataAccessProver) calculateWatingTime() (time.Duration, int64) {
	challengeCycleSeconds := int64((p.interval + p.period).Seconds())
	now := time.Now().Unix()
	duration := now - p.last
	over := duration % challengeCycleSeconds
	var waitingSeconds int64 = 0
	if over < int64(p.interval.Seconds()) {
		waitingSeconds = int64(p.interval.Seconds()) - over
	}

	p.last = p.last + duration - over
	next := p.last + challengeCycleSeconds

	return time.Duration(waitingSeconds) * time.Second, next
}

func (p *DataAccessProver) resetChallengeStatus() error {
	info, err := p.proofInstance.GetChallengeInfo()
	if err != nil {
		return err
	}

	if info.ChalStatus != 0 {
		return p.proofInstance.EndChallenge()
	}

	return nil
}

func (p *DataAccessProver) generateRND() (fr.Element, bool, error) {
	rnd := fr.Element{}
	err := p.proofInstance.GenerateRnd()
	if err != nil {
		return rnd, false, err
	}

	rnd, lock, _, err := p.proofInstance.GetVerifyInfo()
	if err != nil {
		return rnd, false, err
	}
	return rnd, lock, nil
}

func (p *DataAccessProver) selectFiles(rnd fr.Element) ([]bls12381.G1Affine, []kzg.OpeningProof, error) {
	var commits []bls12381.G1Affine = make([]bls12381.G1Affine, p.selectedFileNumber)
	var proofs []kzg.OpeningProof = make([]kzg.OpeningProof, p.selectedFileNumber)
	info, err := p.proofInstance.GetChallengeInfo()
	if err != nil {
		return nil, nil, err
	}
	length := info.ChalLength.Int64()

	rndBytes, err := p.proofInstance.GetRndRawBytes()
	if err != nil {
		return nil, nil, err
	}

	var random *big.Int = big.NewInt(0).SetBytes(rndBytes[:])
	random = new(big.Int).Mod(random, big.NewInt(length))
	startIndex := new(big.Int).Div(random, big.NewInt(2)).Int64()

	var endIndex int64
	if p.selectedFileNumber > length {
		endIndex = startIndex + (length-1)/2
	} else {
		endIndex = startIndex + (p.selectedFileNumber-1)/2
	}

	files, err := GetRangeDAFileInfo(uint(startIndex+1), uint(endIndex+1))
	if err != nil {
		return nil, nil, err
	}

	var tmpCommits = make([]bls12381.G1Affine, len(files))
	var tmpProofs = make([]kzg.OpeningProof, len(files))
	for index, file := range files {
		if file.Expiration > p.last {
			var w bytes.Buffer
			err = daStore.GetObject(context.TODO(), file.Mid, &w, api.ObjectOptions{})
			if err != nil {
				return nil, nil, err
			}

			poly := split(w.Bytes())
			proof, err := kzg.Open(poly, rnd, p.provingKey)
			if err != nil {
				return nil, nil, err
			}

			tmpCommits[index] = file.Commit
			tmpProofs[index] = proof
		} else {
			tmpCommits[index] = zeroCommit
			tmpProofs[index] = zeroProof
		}
	}

	for index := 0; index < int(p.selectedFileNumber); index++ {
		commits[index] = tmpCommits[index%int(length)/2]
		proofs[index] = tmpProofs[index%int(length)/2]
	}

	return commits, proofs, nil
}

func (p *DataAccessProver) proveToContract(commits []bls12381.G1Affine, proofs []kzg.OpeningProof, rnd fr.Element) error {
	// fold proof
	var foldedCommit bls12381.G1Affine
	var foldedProof kzg.OpeningProof
	var foldedPi bls12381.G1Affine
	var foldedValue fr.Element

	foldedCommit = commits[0]
	foldedPi = proofs[0].H
	foldedValue = proofs[0].ClaimedValue
	for index := 1; index < len(commits); index++ {
		// compute
		foldedCommit.Add(&foldedCommit, &commits[index])
		// compute
		foldedPi.Add(&foldedPi, &proofs[index].H)
		// compute
		foldedValue.Add(&foldedValue, &proofs[index].ClaimedValue)
	}

	// var tmpCommit []bls12381.G1Affine
	// var aggregatedCommits [10]bls12381.G1Affine
	// var splitLength = len(commits) / 10
	// for i := 0; i < 10; i++ {
	// 	tmpCommit = commits[i*splitLength : (i+1)*splitLength]
	// 	var aggregatedCommit bls12381.G1Affine = tmpCommit[0]
	// 	for _, commit := range tmpCommit[1:] {
	// 		aggregatedCommit.Add(&aggregatedCommit, &commit)
	// 	}
	// 	aggregatedCommits[i] = aggregatedCommit
	// }

	foldedProof.H = foldedPi
	foldedProof.ClaimedValue = foldedValue

	// err := kzg.Verify(&foldedCommit, &foldedProof, rnd, DefaultSRS.Vk)
	// if err != nil {
	// 	log.Println("verify kzg proof failed", err.Error())
	// }

	// g2, err := p.proofInstance.GetVK()
	// if err != nil {
	// 	log.Println("got vk failed", err.Error())
	// } else {
	// 	if !g2.Equal(&DefaultSRS.Vk.G2[1]) {
	// 		log.Println("vk is not equal")
	// 	}
	// }

	return p.proofInstance.SubmitAggregationProof(rnd, foldedCommit, foldedProof)
}

func (p *DataAccessProver) responseChallenge(commits []bls12381.G1Affine) error {
	var splitedCommits [10][]bls12381.G1Affine
	for {
		info, err := p.proofInstance.GetChallengeInfo()
		if err != nil {
			return err
		}

		if info.ChalStatus%2 == 0 {
			if time.Now().Unix() > p.last+p.respondTime*int64(info.ChalStatus+1) {
				if info.ChalStatus != 0 {
					return p.proofInstance.EndChallenge()
				} else {
					return nil
				}
			}
		} else if info.ChalStatus == 11 {
			var splitLength = len(commits) / 10
			selectedCommits := commits[int(info.ChalIndex)*splitLength : int(info.ChalIndex+1)*splitLength]
			return p.proofInstance.OneStepProve(selectedCommits)
		} else {
			if info.ChalStatus != 1 {
				commits = splitedCommits[info.ChalIndex]
			}
			var aggregatedCommits [10]bls12381.G1Affine
			var splitLength = len(commits) / 10
			for i := 0; i < 10; i++ {
				splitedCommits[i] = commits[i*splitLength : (i+1)*splitLength]
				var aggregatedCommit bls12381.G1Affine = splitedCommits[i][0]
				for _, commit := range splitedCommits[i][1:] {
					aggregatedCommit.Add(&aggregatedCommit, &commit)
				}
				aggregatedCommits[i] = aggregatedCommit
			}
			err := p.proofInstance.ResponseChallenge(aggregatedCommits)
			if err != nil {
				return err
			}
			// data, err := json.MarshalIndent(info, "", "\t")
			// if err != nil {
			// 	return err
			// }
			// log.Println(string(data))
			// log.Println(commitsCopy[info.StartIndex.Int64()%p.selectedFileNumber])
			// commit, err := p.proofInstance.GetSelectFileCommit(info.StartIndex)
			// if err != nil {
			// 	return err
			// }
			// log.Println(commit)
		}
		time.Sleep(5 * time.Second)
	}
}

const ShardingLen = 127

func Pad127(in []byte, res []fr.Element) {
	if len(in) != 127 {
		if len(in) > 127 {
			in = in[:127]
		} else {
			padding := make([]byte, 127-len(in))
			in = append(in, padding...)
		}
	}

	tmp := make([]byte, 32)
	copy(tmp[:31], in[:31])

	t := in[31] >> 6
	tmp[31] = in[31] & 0x3f
	res[0].SetBytes(tmp)

	var v byte
	for i := 32; i < 64; i++ {
		v = in[i]
		tmp[i-32] = (v << 2) | t
		t = v >> 6
	}
	t = v >> 4
	tmp[31] &= 0x3f
	res[1].SetBytes(tmp)

	for i := 64; i < 96; i++ {
		v = in[i]
		tmp[i-64] = (v << 4) | t
		t = v >> 4
	}
	t = v >> 2
	tmp[31] &= 0x3f
	res[2].SetBytes(tmp)

	for i := 96; i < 127; i++ {
		v = in[i]
		tmp[i-96] = (v << 6) | t
		t = v >> 2
	}
	tmp[31] = t & 0x3f
	res[3].SetBytes(tmp)
}

func split(data []byte) []fr.Element {
	num := (len(data)-1)/ShardingLen + 1

	atom := make([]fr.Element, num*4)

	padding := make([]byte, ShardingLen*num-len(data))
	data = append(data, padding...)

	for i := 0; i < num; i++ {
		Pad127(data[ShardingLen*i:ShardingLen*(i+1)], atom[4*i:4*i+4])
	}

	return atom
}
