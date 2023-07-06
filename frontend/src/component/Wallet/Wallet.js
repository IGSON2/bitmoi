import { MetaMaskSDK } from "@metamask/sdk";
import { useEffect, useState } from "react";
import Web3 from "web3";
import ContractABI from "../../contract/moiABI.json";

function Wallet() {
  const [hasProvider, setHasProvider] = useState(false);
  const [wallet, setWallet] = useState({ accounts: [] });
  const [contractInstance, setContractInstance] = useState(null);

  useEffect(() => {
    if (window.ethereum) {
      setHasProvider(true);
      const web3 = new Web3(window.ethereum);
      const contract = new web3.eth.Contract(
        ContractABI,
        "0x8845316dd44894FFFcAc0F6184aa0b4E5a52e1D9"
      );
      setContractInstance(contract);
    } else {
      setHasProvider(false);
      alert("install metamask extension!!");
    }
  }, []);
  const updateWallet = async (accounts) => {
    setWallet({ accounts });
  };

  console.log(contractInstance.methods.balanceOf, wallet.accounts[0]);

  const handleConnect = async () => {
    let accounts = await window.ethereum.request({
      method: "eth_requestAccounts",
    });
    updateWallet(accounts);
  };

  const getTokenBalance = async () => {
    try {
      const result = await contractInstance.methods
        .balanceOf(wallet.accounts[0])
        .call(); // getter는 call
      console.log(result);
    } catch (error) {
      console.error(error);
    }
  };

  const connectWallet = () => {};
  const getAccount = () => {};
  return (
    <div>
      <button onClick={connectWallet}>connect</button>
      <button onClick={getAccount}>getProvider</button>
      <div>Injected Provider {hasProvider ? "DOES" : "DOES NOT"} Exist</div>
      {hasProvider && <button onClick={handleConnect}>Connect MetaMask</button>}
      {wallet.accounts.length > 0 && (
        <div>Wallet Accounts: {wallet.accounts[0]}</div>
      )}
      <button onClick={getTokenBalance}>getTokenBalacne</button>
    </div>
  );
}

export default Wallet;
