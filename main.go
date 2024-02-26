package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var SanctionList = map[string]struct{}{
	"0x8589427373d6d84e98730d7795d8f6f8731fda16": {},
	"0x722122df12d4e14e13ac3b6895a86e84145b6967": {},
	"0xdd4c48c0b24039969fc16d1cdf626eab821d3384": {},
	"0xd90e2f925da726b50c4ed8d0fb90ad053324f31b": {},
	"0xd96f2b1c14db8458374d9aca76e26c3d18364307": {},
	"0x4736dcf1b7a3d580672cce6e7c65cd5cc9cfba9d": {},
	"0xd4b88df4d29f5cedd6857912842cff3b20c8cfa3": {},
	"0x910cbd523d972eb0a6f4cae4618ad62622b39dbf": {},
	"0xa160cdab225685da1d56aa342ad8841c3b53f291": {},
	"0xfd8610d20aa15b7b2e3be39b396a1bc3516c7144": {},
	"0xf60dd140cff0706bae9cd734ac3ae76ad9ebc32a": {},
	"0x22aaa7720ddd5388a3c0a3333430953c68f1849b": {},
	"0xba214c1c1928a32bffe790263e38b4af9bfcd659": {},
	"0xb1c8094b234dce6e03f10a5b673c1d8c69739a00": {},
	"0x527653ea119f3e6a1f5bd18fbf4714081d7b31ce": {},
	"0x58e8dcc13be9780fc42e8723d8ead4cf46943df2": {},
	"0xd691f27f38b395864ea86cfc7253969b409c362d": {},
	"0xaeaac358560e11f52454d997aaff2c5731b6f8a6": {},
	"0x1356c899d8c9467c7f71c195612f8a395abf2f0a": {},
	"0xa60c772958a3ed56c1f15dd055ba37ac8e523a0d": {},
	"0x169ad27a470d064dede56a2d3ff727986b15d52b": {},
	"0x0836222f2b2b24a3f36f98668ed8f0b38d1a872f": {},
	"0xf67721a2d8f736e75a49fdd7fad2e31d8676542a": {},
	"0x9ad122c22b14202b4490edaf288fdb3c7cb3ff5e": {},
	"0x905b63fff465b9ffbf41dea908ceb12478ec7601": {},
	"0x07687e702b410fa43f4cb4af7fa097918ffd2730": {},
	"0x94a1b5cdb22c43faab4abeb5c74999895464ddaf": {},
	"0xb541fc07bc7619fd4062a54d96268525cbc6ffef": {},
	"0x12d66f87a04a9e220743712ce6d9bb1b5616b8fc": {},
	"0x47ce0c6ed5b0ce3d3a51fdb1c52dc66a7c3c2936": {},
	"0x23773e65ed146a459791799d01336db287f25334": {},
	"0xd21be7248e0197ee08e0c20d4a96debdac3d20af": {},
	"0x610b717796ad172b316836ac95a2ffad065ceab4": {},
	"0x178169b423a011fff22b9e3f3abea13414ddd0f1": {},
	"0xbb93e510bbcd0b7beb5a853875f9ec60275cf498": {},
	"0x2717c5e28cf931547b621a5dddb772ab6a35b701": {},
	"0x03893a7c7463ae47d46bc7f091665f1893656003": {},
	"0xca0840578f57fe71599d29375e16783424023357": {},
	"0x3aac1cc67c2ec5db4ea850957b967ba153ad6279": {},
	"0x76d85b4c0fc497eecc38902397ac608000a06607": {},
	"0x0e3a09dda6b20afbb34ac7cd4a6881493f3e7bf7": {},
	"0x723b78e67497e85279cb204544566f4dc5d2aca0": {},
	"0xcc84179ffd19a1627e79f8648d09e095252bc418": {},
	"0x6bf694a291df3fec1f7e69701e3ab6c592435ae7": {},
	"0x330bdfade01ee9bf63c209ee33102dd334618e0a": {},
	"0xa5c2254e4253490c54cef0a4347fddb8f75a4998": {},
	"0xaf4c0b70b2ea9fb7487c7cbb37ada259579fe040": {},
	"0xdf231d99ff8b6c6cbf4e9b9a945cbacef9339178": {},
	"0x1e34a77868e19a6647b1f2f47b51ed72dede95dd": {},
	"0xd47438c816c9e7f2e2888e060936a499af9582b3": {},
	"0x84443cfd09a48af6ef360c6976c5392ac5023a1f": {},
	"0xd5d6f8d9e784d0e26222ad3834500801a68d027d": {},
	"0xaf8d1839c3c67cf571aa74b5c12398d4901147b3": {},
	"0x407cceeaa7c95d2fe2250bf9f2c105aa7aafb512": {},
	"0x05e0b5b40b7b66098c2161a5ee11c5740a3a7c45": {},
	"0xd8d7de3349ccaa0fde6298fe6d7b7d0d34586193": {},
	"0x3efa30704d2b8bbac821307230376556cf8cc39e": {},
	"0x746aebc06d2ae31b71ac51429a19d54e797878e9": {},
	"0x5f6c97c6ad7bdd0ae7e0dd4ca33a4ed3fdabd4d7": {},
	"0xf4b067dd14e95bab89be928c07cb22e3c94e0daa": {},
	"0x01e2919679362dfbc9ee1644ba9c6da6d6245bb1": {},
	"0x2fc93484614a34f26f7970cbb94615ba109bb4bf": {},
	"0x26903a5a198d571422b2b4ea08b56a37cbd68c89": {},
	"0xb20c66c4de72433f3ce747b58b86830c459ca911": {},
	"0x2573bac39ebe2901b4389cd468f2872cf7767faf": {},
	"0x653477c392c16b0765603074f157314cc4f40c32": {},
	"0x88fd245fedec4a936e700f9173454d1931b4c307": {},
	"0x09193888b3f38c82dedfda55259a82c0e7de875e": {},
	"0x5cab7692d4e94096462119ab7bf57319726eed2a": {},
	"0x756c4628e57f7e7f8a459ec2752968360cf4d1aa": {},
	"0xd82ed8786d7c69dc7e052f7a542ab047971e73d2": {},
	"0x77777feddddffc19ff86db637967013e6c6a116c": {},
	"0x833481186f16cece3f1eeea1a694c42034c3a0db": {},
	"0xb04e030140b30c27bcdfaafffa98c57d80eda7b4": {},
	"0xcee71753c9820f063b38fdbe4cfdaf1d3d928a80": {},
	"0x8281aa6795ade17c8973e1aedca380258bc124f9": {},
	"0x57b2b8c82f065de8ef5573f9730fc1449b403c9f": {},
	"0x23173fe8b96a4ad8d2e17fb83ea5dcccdca1ae52": {},
	"0x538ab61e8a9fc1b2f93b3dd9011d662d89be6fe6": {},
	"0x94be88213a387e992dd87de56950a9aef34b9448": {},
	"0x242654336ca2205714071898f67e254eb49acdce": {},
	"0x776198ccf446dfa168347089d7338879273172cf": {},
	"0xedc5d01286f99a066559f60a585406f3878a033e": {},
	"0xd692fd2d0b2fbd2e52cfa5b5b9424bc981c30696": {},
	"0xdf3a408c53e5078af6e8fb2a85088d46ee09a61b": {},
	"0x743494b60097a2230018079c02fe21a7b687eaa5": {},
	"0x94c92f096437ab9958fc0a37f09348f30389ae79": {},
	"0x5efda50f22d34f262c29268506c5fa42cb56a1ce": {},
	"0x2f50508a8a3d323b91336fa3ea6ae50e55f32185": {},
	"0x179f48c78f57a3a78f0608cc9197b8972921d1d2": {},
	"0xffbac21a641dcfe4552920138d90f3638b3c9fba": {},
}

var fileAddresses = map[string]struct{}{}

var (
	nodeURL  = flag.String("node-url", "", "Ethereum node URL")
	fileName = flag.String("file-name", "addresses.json", "path to compliance list")
)

func main() {
	flag.Parse()

	b, err := os.ReadFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(b, &fileAddresses); err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(*nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	contracts := []common.Address{}
	for k := range SanctionList {
		contracts = append(contracts, common.HexToAddress(k))
	}

	// filter for logs from all tornado cash contracts
	query := ethereum.FilterQuery{
		Addresses: contracts,
	}

	logs := make(chan types.Log)

	// subscribe to the logs
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	// load the tornado cash contract ABI
	abiBytes, err := os.ReadFile("abi.json")
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(bytes.NewReader(abiBytes))
	if err != nil {
		log.Fatal(err)
	}

	writeToFileTicker := time.NewTicker(5 * time.Minute)

	// TODO: need reconnection logic
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			_, err := contractAbi.Unpack("Deposit", vLog.Data)
			if err != nil {
				continue
			}

			tx, _, err := client.TransactionByHash(context.Background(), vLog.TxHash)
			if err != nil {
				fmt.Println(err)
				continue
			}

			from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fileAddresses[from.String()] = struct{}{}
		case <-writeToFileTicker.C:

			fileBytes, err := json.Marshal(fileAddresses)
			if err != nil {
				fmt.Println("could not write to file", err)
				continue
			}
			// write to file
			if err := os.WriteFile(*fileName, fileBytes, 0644); err != nil {
				fmt.Println("could not write to file", err)
				continue
			}
		}
	}

}
