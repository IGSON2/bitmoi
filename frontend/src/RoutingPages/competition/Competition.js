import TradingBoard from "../../component/tradingBoard/TradingBoard";
import styles from "./competition.module.css";

function Competition() {
  const scoreId = Date.now().toString();
  return (
    <div className={styles.competitionpage}>
      <div className={styles.chart}>
        <TradingBoard
          modeHeight={0.83}
          mode={"competition"}
          scoreId={scoreId}
        />
      </div>
    </div>
  );
}
export default Competition;
