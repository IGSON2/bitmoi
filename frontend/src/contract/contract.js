import Web3 from "web3";
import ContractABI from "./moiABI.json";

const getBalance = async () => {
  const web3 = new Web3(window.ethereum);
  const contract = new web3.eth.Contract(
    ContractABI,
    "0x3CF4A2Deaf37eA7B7e2413C6ce8C35215099cDc7"
  );

  let accounts = await window.ethereum.request({
    method: "eth_requestAccounts",
  });

  let balance = await contract.methods.balanceOf(accounts[0]).call();

  return Number(balance);
};

const getAccount = async () => {
  let accounts = await window.ethereum.request({
    method: "eth_requestAccounts",
  });

  return accounts[0];
};

const getChainId = async () => {
  const web3 = new Web3(window.ethereum);
  const currentChainId = await web3.eth.net.getId();
  return currentChainId;
};

export { getBalance, getAccount, getChainId };
