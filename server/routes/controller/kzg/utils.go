package kzg

func Split(data []byte) []Fr {
	num := (len(data)-1)/ShardingLen + 1

	atom := make([]Fr, num)

	for i := 0; i < num-1; i++ {
		atom[i].SetBytes(data[ShardingLen*i : ShardingLen*(i+1)])
	}

	atom[num-1].SetBytes(data[ShardingLen*(num-1):])

	return atom
}
