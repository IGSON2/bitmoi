import Timeformatter from "../../../component/Timeformatter/Timeformatter";
import styles from "./ScoreGraph.module.css";

function ScoreGraph({ obj, index }) {
  const formattedID = Timeformatter(obj.scoreid, false);
  return (
    <div className={styles.wholegraph}>
      <div className={styles.fieldshort}>{index}</div>
      <div className={styles.fieldlong}>{formattedID}</div>
      <div
        className={styles.fieldlong}
        style={
          obj.pairname.includes("(Liq.)")
            ? { color: "#ef5350" }
            : { color: "black" }
        }
      >
        {obj.pairname}
      </div>
      <div className={styles.fieldlong}>{obj.entrytime}</div>
      <div className={styles.fieldshort}>{obj.stage}</div>
      <div className={styles.field}>{obj.position}</div>
      <div className={styles.fieldshort}>X{obj.leverage}</div>
      <div className={styles.field}>{obj.entryprice}</div>
      <div className={styles.field}>{obj.endprice}</div>
      <div className={styles.fieldshort}>{obj.outtime}H</div>
      <div className={styles.field}>{Math.floor(obj.pnl * 100) / 100}</div>
      <div className={styles.field}>{obj.roe} %</div>
    </div>
  );
}

export default ScoreGraph;
