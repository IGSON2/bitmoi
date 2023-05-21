import TradingBoard from "../../component/tradingBoard/TradingBoard";
import styles from "./competition.module.css";

function Competition() {
  return (
    <div className={styles.competitionpage}>
      <div className={styles.chart}>
        <TradingBoard modeHeight={0.83} mode={"competition"} />
      </div>
    </div>
  );
}
export default Competition;
