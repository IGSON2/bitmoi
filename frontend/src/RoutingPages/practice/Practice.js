import TradingBoard from "../../component/tradingBoard/TradingBoard";
import mockupimg from "../../component/images/mockup_prac.png";
import styles from "./practice.module.css";
import { useState } from "react";

function Practice() {
  const score_id = Date.now().toString();
  const [isLoaded, setIsLoaded] = useState(false);

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
        {isLoaded ? <img src={mockupimg}></img> : null}
      </div>
    </div>
  );
}
export default Practice;
