import TradingBoard from "../../component/tradingBoard/TradingBoard";
import mockupimg from "../../component/images/mockupad.png";
import styles from "./practice.module.css";
import { useState } from "react";

function Practice() {
  return (
    <div className={styles.practicepage}>
      <div className={styles.chart}>
        <TradingBoard modeHeight={0.75} mode={"practice"} />
      </div>
      <div className={styles.ad}>
        <img src={mockupimg}></img>
      </div>
    </div>
  );
}
export default Practice;
