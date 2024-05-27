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
  const [disableReq, setDisAbleReq] = useState(false);
  const getFreeToken = async () => {
    try {
      const res = await axiosClient.post("/auth/freeToken", {
        addr: addr,
      });
      setWarning("");
    } catch (error) {
      const resMsg = error.response.data;
      if (error.response.status === 429) {
        setDisAbleReq(true);
        setWarning("요청이 너무 잦습니다.");
        return;
      }
      if (resMsg.includes("allowance")) {
        if (resMsg[7] === ".") {
          setWarning(`다음 무료 발급까지 남은 시간 [ ${resMsg.slice(0, 7)}s ]`);
        } else {
          setWarning(`다음 무료 발급까지 남은 시간 [ ${resMsg.slice(0, 8)}s ]`);
        }
      }
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
        <button
          className={styles.faucet}
          onClick={getFreeToken}
          disabled={disableReq}
        >
          토큰 발급받기
        </button>
        {warning ? <div className={styles.warning}>{warning}</div> : null}
      </div>
    </div>
  );
}

export default Freetoken;
