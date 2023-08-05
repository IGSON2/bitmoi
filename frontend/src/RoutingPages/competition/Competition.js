import TradingBoard from "../../component/tradingBoard/TradingBoard";
import styles from "./competition.module.css";

function Competition() {
  const score_id = Date.now().toString();
  return (
    <div className={styles.competitionpage}>
      <div className={styles.chart}>
        <TradingBoard
          modeHeight={0.83}
          mode={"competition"}
          score_id={score_id}
        />
      </div>
    </div>
  );
}
export default Competition;
