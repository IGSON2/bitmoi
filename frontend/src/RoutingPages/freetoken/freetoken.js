import axiosClient from "../../component/backendConn/axiosClient";
import { useEffect, useState } from "react";
import { getAccount } from "../../contract/contract";
import styles from "./freetoken.module.css";
import H_NavBar from "../../component/navbar/H_NavBar";

function Freetoken() {
  const [addr, setAddr] = useState("");
  const getFreeToken = async () => {
    await axiosClient.post("/freeToken", {
      addr: addr,
    });
  };

  useEffect(() => {
    const initAddr = async () => {
      try {
        const resAddr = await getAccount();
        if (resAddr !== "") {
          setAddr(addr);
        }
      } catch (error) {
        console.error(error);
      }
    };
    initAddr();
  }, []);

  return (
    <div className={styles.wrapper}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.mainbody}>
        <button className={styles.faucet} onClick={getFreeToken}>
          토큰 발급받기
        </button>
      </div>
    </div>
  );
}

export default Freetoken;
