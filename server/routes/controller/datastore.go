package controller

import (
	"context"
	"fmt"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) getSpaceCheckInfo(ctx context.Context, address string, size uint64) (api.CheckInfo, error) {
	var res api.CheckInfo

	pi, err := c.SpacePayInfo(ctx, address)
	if err != nil {
		return res, err
	}

	cs, err := c.datastore.GetSpaceInfo(ctx, address)
	if err != nil {
		return res, err
	}
	checksize := cs.FileSize.Uint64() + size

	cs.FileSize.SetUint64(checksize)
	cs.Nonce = pi.Nonce
	if checksize > pi.FreeByte+pi.SizeByte {
		lerr := logs.ControllerError{Message: fmt.Sprintf("space not enough, have %d, need %d", pi.FreeByte+pi.SizeByte, checksize)}
		logger.Error(lerr)
		return res, lerr
	}

	return cs, nil
}
func (c *Controller) getTrafficCheckInfo(ctx context.Context, address string, size uint64) (api.CheckInfo, error) {
	var res api.CheckInfo

	pi, err := c.TrafficPayInfo(ctx, address)
	if err != nil {
		return res, err
	}

	cs, err := c.datastore.GetTrafficInfo(ctx, address)
	if err != nil {
		return res, err
	}
	checksize := cs.FileSize.Uint64() + size

	cs.FileSize.SetUint64(checksize)
	cs.Nonce = pi.Nonce
	if checksize > pi.FreeByte+pi.SizeByte {
		lerr := logs.ControllerError{Message: fmt.Sprintf("space not enough, have %d, need %d", pi.FreeByte+pi.SizeByte, checksize)}
		logger.Error(lerr)
		return res, lerr
	}

	return cs, nil
}
