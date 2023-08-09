import TradingBoard from "../../component/tradingBoard/TradingBoard";
import mockupimg from "../../component/images/mockupad.png";
import styles from "./practice.module.css";

function Practice() {
  const score_id = Date.now().toString();

  return (
    <div className={styles.practicepage}>
      <div className={styles.chart}>
        <TradingBoard modeHeight={0.75} mode={"practice"} score_id={score_id} />
      </div>
      <div className={styles.ad}>
        <img src={mockupimg}></img>
      </div>
    </div>
  );
}
export default Practice;
