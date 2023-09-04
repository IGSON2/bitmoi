import TradingBoard from "../../component/tradingBoard/TradingBoard";
import mockupimg from "../../component/images/mockup_prac.png";
import styles from "./practice.module.css";
import { useEffect, useState } from "react";
import getSelectedBidderImg from "../../component/backendConn/getSelectedBidderImg";

function Practice() {
  const score_id = Date.now().toString();
  const [isLoaded, setIsLoaded] = useState(false);

  const [imgLink, setImgLink] = useState("");

  const adClick = () => {
    window.open("/ad-bidding/practice", "_blank");
  };

  useEffect(() => {
    const getBidder = async () => {
      const img = await getSelectedBidderImg("practice");
      img === "" ? setImgLink(mockupimg) : setImgLink(img);
    };
    getBidder();
  }, {});

  console.log(imgLink);
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
        {isLoaded ? <img src={imgLink} onClick={adClick}></img> : null}
      </div>
    </div>
  );
}
export default Practice;
