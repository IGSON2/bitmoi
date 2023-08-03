import { useEffect, useState } from "react";
import Web3 from "web3";
import styles from "./Wallet.module.css";
import ContractABI from "../../contract/moiABI.json";
import moilogo from "../images/logosmall.png";

function Wallet() {
  const baobabTestNetID = 1001;

  const [buttonDisabled, setButtonDisabled] = useState(false);
  const [hasProvider, setHasProvider] = useState(false);
  const [wallet, setWallet] = useState({ accounts: [] });
  const [contractInstance, setContractInstance] = useState(null);
  const [tokenBalance, setTokenBalance] = useState(0);

  const callBalanceOf = async (accountAddress) => {
    await contractInstance.methods
      .balanceOf(accountAddress)
      .call()
      .then((balance) => {
        setTokenBalance(Number(balance));
      })
      .catch((error) => {
        console.error("Error calling balanceOf:", error);
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
      console.log(chainId, `0x${baobabTestNetID.toString(16)}`);
      setButtonDisabled(true);
      setTokenBalance(0);
      alert("only can use token at baobab network.");
    } else {
      setButtonDisabled(false);
      callBalanceOf(wallet.accounts[0]);
    }
    console.log("Chain changed:", chainId);
  };

  const handleAccountsChange = async (accounts) => {
    updateWallet(accounts);
    callBalanceOf(accounts[0]);
    console.log("Accounts changed:", accounts);
  };

  useEffect(() => {
    const initwallet = async () => {
      if (window.ethereum) {
        setHasProvider(true);
        await handleConnect();
        const web3 = new Web3(window.ethereum);
        const currentChainId = await web3.eth.net.getId();
        if (Number(currentChainId) !== baobabTestNetID) {
          const res = await window.ethereum.request({
            method: "wallet_switchEthereumChain",
            params: [{ chainId: `0x${baobabTestNetID.toString(16)}` }],
          });
        } else {
          console.log("Already connected to chain ID:", baobabTestNetID);
        }

        const contract = new web3.eth.Contract(
          ContractABI,
          "0xf4CFFdF8032B7C59d8254538Cc9F3f20BF2a03fF"
        );
        setContractInstance(contract);

        window.ethereum.on("chainChanged", handleChainChange);
        window.ethereum.on("accountsChanged", handleAccountsChange);
      } else {
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
      {hasProvider ? (
        <div className={styles.balance}>
          <div
            className={styles.logo}
            title={wallet.accounts.length > 0 ? wallet.accounts[0] : ""}
          >
            <img src={moilogo}></img>
          </div>
          <div className={styles.number}>{tokenBalance} MOI</div>
        </div>
      ) : (
        <div className={styles.warning}>Metamask를 설치해 주세요!</div>
      )}
    </div>
  );
}

export default Wallet;
