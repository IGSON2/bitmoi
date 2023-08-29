import TradingBoard from "../../component/tradingBoard/TradingBoard";
import mockupimg from "../../component/images/mockup_prac.png";
import styles from "./practice.module.css";
import { useEffect, useState } from "react";
import getSelectedBidder from "../../component/backendConn/getSelectedBidder";

function Practice() {
  const score_id = Date.now().toString();
  const [isLoaded, setIsLoaded] = useState(false);

  const [imgLink, setImgLink] = useState("");

  const adClick = () => {
    window.open("/ad-bidding/practice", "_blank");
  };

  useEffect(() => {
    //TODO

    getSelectedBidder("practice");
    setImgLink();
  }, {});

  return (
    <div className={styles.practicepage}>
      <div className={styles.chart}>
        <TradingBoard
          modeHeight={0.75}
          mode={"practice"}
          score_id={score_id}
          setIsLoaded={setIsLoaded}
        />
      </div>
      <div className={styles.ad}>
        {isLoaded ? <img src={mockupimg} onClick={adClick}></img> : null}
      </div>
    </div>
  );
}
export default Practice;
