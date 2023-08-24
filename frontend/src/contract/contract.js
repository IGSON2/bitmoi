import Web3 from "web3";
import ContractABI from "./moiABI.json";

const getBalance = async () => {
  const web3 = new Web3(window.ethereum);
  const contract = new web3.eth.Contract(
    ContractABI,
    "0xd7003A3DE013054B52978b02FaE6620f3f9C2368"
  );

  let accounts = await window.ethereum.request({
    method: "eth_requestAccounts",
  });

  let balance = await contract.methods.balanceOf(accounts[0]).call();

  return Number(balance);
};

export default getBalance;
