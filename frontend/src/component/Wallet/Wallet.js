import { useEffect, useState } from "react";
import Web3 from "web3";
import styles from "./Wallet.module.css";
import ContractABI from "../../contract/moiABI.json";
import moilogo from "../images/logo.png";

function Wallet() {
  const baobabTestNetID = 1001; // 8217;

  const [hasProvider, setHasProvider] = useState(false);
  const [wallet, setWallet] = useState({ accounts: [] });
  const [contractInstance, setContractInstance] = useState(null);
  const [tokenBalance, setTokenBalance] = useState(0);
  const [warning, setWarning] = useState("");
  const [needsToInstall, setNeedsToInstall] = useState(false);
  const [needsToConnGuide, setNeedsToConnGuide] = useState(false);

  const callBalanceOf = async (accountAddress) => {
    await contractInstance.methods
      .balanceOf(accountAddress)
      .call()
      .then((balance) => {
        setTokenBalance(Number(balance));
        setWarning("");
      })
      .catch((error) => {
        console.error("Error calling balanceOf:", error);
        setWarning("Klaytn 메인넷에 연결되어 있는지 확인해주세요.");
      });
  };

  const updateWallet = (accounts) => {
    setWallet({ accounts });
  };

  const handleConnect = async () => {
    let accounts = await window.ethereum.request({
      method: "eth_requestAccounts",
    });
    updateWallet(accounts);
  };

  const handleChainChange = (chainId) => {
    if (chainId !== `0x${baobabTestNetID.toString(16)}`) {
      setTokenBalance(0);
      setWarning("Metamask에 Klaytn 메인넷을 등록해 주세요");
      setNeedsToConnGuide(true);
    } else {
      callBalanceOf(wallet.accounts[0]);
      setWarning("");
      setNeedsToConnGuide(false);
    }
    setNeedsToInstall(false);
  };

  const handleAccountsChange = async (accounts) => {
    updateWallet(accounts);
    callBalanceOf(accounts[0]);
  };

  useEffect(() => {
    const initwallet = async () => {
      if (window.ethereum) {
        setWarning("Metamask에 로그인되어 있지 않습니다.");
        await handleConnect();
        setWarning("");
        const web3 = new Web3(window.ethereum);
        const currentChainId = await web3.eth.net.getId();
        if (Number(currentChainId) !== baobabTestNetID) {
          try {
            await window.ethereum.request({
              method: "wallet_switchEthereumChain",
              params: [{ chainId: `0x${baobabTestNetID.toString(16)}` }],
            });
            setNeedsToConnGuide(false);
          } catch (error) {
            if (error.code === 4902) {
              setWarning("Metamask에 Klaytn 메인넷을 등록해 주세요");
              setNeedsToConnGuide(true);
            }
          }
        }

        const contract = new web3.eth.Contract(
          ContractABI,
          "0xC6D6e6E081dA15bd4258bDbdC121AD2D72ad793B"
        );
        setContractInstance(contract);

        window.ethereum.on("chainChanged", handleChainChange);
        window.ethereum.on("accountsChanged", handleAccountsChange);
        setHasProvider(true);
        setNeedsToInstall(false);
      } else {
        setWarning("Metamask가 설치되어 있지 않습니다.");
        setNeedsToInstall(true);
        setHasProvider(false);
      }
    };

    initwallet();

    return () => {
      if (window.ethereum) {
        window.ethereum.removeListener("chainChanged", handleChainChange);
        window.ethereum.removeListener("accountsChanged", handleAccountsChange);
      }
    };
  }, []);

  useEffect(() => {
    if (wallet.accounts[0]) {
      callBalanceOf(wallet.accounts[0]);
    }
  }, [contractInstance]);

  return (
    <div className={styles.wallet}>
      {hasProvider && warning === "" ? (
        <div className={styles.balance}>
          <div
            className={styles.logo}
            title={wallet.accounts.length > 0 ? wallet.accounts[0] : ""}
          >
            <img src={moilogo}></img>
          </div>
          <div
            className={styles.number}
            title={tokenBalance.toLocaleString("ko-KR")}
          >
            {String(tokenBalance).length > 4
              ? String(tokenBalance).slice(0, 4) + ".."
              : tokenBalance}
            MOI
          </div>
        </div>
      ) : (
        <div className={styles.warning}>
          <div className={styles.warningtext}>{warning}</div>
          <div className={styles.linkbox}>
            {needsToInstall ? (
              <a href="https://metamask.io/" target="_blank">
                설치하기
              </a>
            ) : needsToConnGuide ? (
              <a
                href="https://ko.docs.klaytn.foundation/content/dapp/tutorials/connecting-metamask#install-metamask"
                target="_blank"
              >
                등록방법
              </a>
            ) : null}
          </div>
        </div>
      )}
    </div>
  );
}

export default Wallet;
