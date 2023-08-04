import Web3 from "web3";
import ContractABI from "./moiABI.json";

const getBalance = async () => {
  const web3 = new Web3(window.ethereum);
  const contract = new web3.eth.Contract(
    ContractABI,
    "0xf4CFFdF8032B7C59d8254538Cc9F3f20BF2a03fF"
  );

  let accounts = await window.ethereum.request({
    method: "eth_requestAccounts",
  });

  let balance = await contract.methods.balanceOf(accounts[0]).call();

  return Number(balance);
};

export default getBalance;
