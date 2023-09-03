import axiosClient from "../../component/backendConn/axiosClient";
import { useEffect, useState } from "react";
import { getAccount } from "../../contract/contract";
import styles from "./freetoken.module.css";
import H_NavBar from "../../component/navbar/H_NavBar";
import checkAccessTokenValidity from "../../component/backendConn/checkAccessTokenValidity";

function Freetoken() {
  const [addr, setAddr] = useState("");
  // const [userInfo, setUserInfo] = useState();
  const [warning, setWarning] = useState("");
  const getFreeToken = async () => {
    try {
      const res = await axiosClient.post("/freeToken", {
        addr: addr,
      });
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    const initAddr = async () => {
      const infoRes = await checkAccessTokenValidity();
      if (infoRes.metamask_address !== "") {
        setAddr(infoRes.metamask_address);
      } else {
        try {
          const resAddr = await getAccount();
          if (resAddr !== "") {
            setAddr(resAddr);
          }
        } catch (error) {
          setWarning("Metamask에 등록된 계좌가 없습니다.");
          console.error(error);
        }
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
        {warning ? <div className={styles.warning}>{warning}</div> : null}
      </div>
    </div>
  );
}

export default Freetoken;
