import { useEffect, useState } from "react";
import Web3 from "web3";
import ContractABI from "../../contract/moiABI.json";

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

  const updateWallet = async (accounts) => {
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
    await updateWallet(accounts);
    callBalanceOf(accounts[0]);
    console.log("Accounts changed:", accounts);
  };

  useEffect(() => {
    if (window.ethereum) {
      setHasProvider(true);
      handleConnect();
      const web3 = new Web3(window.ethereum);
      web3.eth.net
        .getId()
        .then((currentChainId) => {
          if (currentChainId !== baobabTestNetID) {
            window.ethereum
              .request({
                method: "wallet_switchEthereumChain",
                params: [{ chainId: `0x${baobabTestNetID.toString(16)}` }],
              })
              .then(() => {
                console.log("Connected to chain ID:", baobabTestNetID);
              })
              .catch((error) => {
                console.error("Failed to switch chain ID:", error);
              });
          } else {
            console.log("Already connected to chain ID:", baobabTestNetID);
          }
        })
        .catch((error) => {
          console.error("Error retrieving current chain ID:", error);
        })
        .then(() => {
          callBalanceOf(wallet.accounts[0]);
          console.log(wallet.accounts[0], tokenBalance);
        });
      const contract = new web3.eth.Contract(
        ContractABI,
        "0xf4CFFdF8032B7C59d8254538Cc9F3f20BF2a03fF"
      );
      setContractInstance(contract);

      window.ethereum.on("chainChanged", handleChainChange);
      window.ethereum.on("accountsChanged", handleAccountsChange);

      return () => {
        if (window.ethereum) {
          window.ethereum.removeListener("chainChanged", handleChainChange);
          window.ethereum.removeListener(
            "accountsChanged",
            handleAccountsChange
          );
        }
      };
    } else {
      setHasProvider(false);
      alert("install metamask extension!!");
    }
  }, []);

  return (
    <div>
      {hasProvider ? (
        <div>
          {wallet.accounts.length > 0 && (
            <div>Wallet Accounts: {wallet.accounts[0]}</div>
          )}
          <div>Token Balance: {tokenBalance}</div>
          <button></button>
        </div>
      ) : (
        <div></div>
      )}
    </div>
  );
}

export default Wallet;
