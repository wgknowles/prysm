package sync

import (
	"context"
	"time"

	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
)

// metaDataHandler reads the incoming metadata rpc request from the peer.
func (r *Service) metaDataHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	defer func() {
		if err := stream.Close(); err != nil {
			log.WithError(err).Error("Failed to close stream")
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	setRPCStreamDeadlines(stream)

	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		return err
	}
	_, err := r.p2p.Encoding().EncodeWithLength(stream, r.p2p.Metadata())
	return err
}

func (r *Service) sendMetaDataRequest(ctx context.Context, id peer.ID) (*pb.MetaData, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stream, err := r.p2p.Send(ctx, new(interface{}), p2p.RPCMetaDataTopic, id)
	if err != nil {
		return nil, err
	}
	// we close the stream outside of `send` because
	// metadata requests send no payload, so closing the
	// stream early leads it to a reset.
	defer func() {
		if err := stream.Close(); err != nil {
			log.WithError(err).Error("Failed to close stream")
		}
	}()
	code, errMsg, err := ReadStatusCode(stream, r.p2p.Encoding())
	if err != nil {
		return nil, err
	}
	if code != 0 {
		r.p2p.Peers().IncrementBadResponses(stream.Conn().RemotePeer())
		return nil, errors.New(errMsg)
	}
	msg := new(pb.MetaData)
	if err := r.p2p.Encoding().DecodeWithLength(stream, msg); err != nil {
		r.p2p.Peers().IncrementBadResponses(stream.Conn().RemotePeer())
		return nil, err
	}
	return msg, nil
}
