package api

import (
	"bitmoi/backend/contract"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *Server) BiddingLoop() error {

	for {
		biddingTimer := time.NewTimer(s.biddingDuration)
		select {
		case <-biddingTimer.C:
			hash, err := s.erc20Contract.UnLockTokens(contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
			if err != nil {
				var initErr error
				s.erc20Contract, initErr = contract.InitErc20Contract(s.config.PrivateKey)
				if initErr != nil {
					log.Err(err).Msgf("Cannot unlock token. stop server..")
					os.Exit(100)
				}
				continue
			}
			receipt, err := s.erc20Contract.WaitAndReturnTxReceipt(hash)
			if err != nil || receipt == nil {
				var initErr error
				s.erc20Contract, initErr = contract.InitErc20Contract(s.config.PrivateKey)
				if initErr != nil {
					log.Err(err).Msgf("Cannot get receipt of unlock token transaction. stop server..")
					os.Exit(100)
				}
				continue
			}
			s.nextUnlockDate = time.Now().Add(s.biddingDuration)
			log.Info().Msgf("Unlock token successfully. hash: %s, next unlock date: %s", hash.Hex(), s.nextUnlockDate.Format("2006-01-02 15:04:05"))
		case <-s.exitCh:
			return nil
		}
	}
}
