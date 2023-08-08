import styles from "./MoreInfo.module.css";

function MoreInfo({ data, obj }) {
  return (
    <div className={styles.moreinfo}>
      <h3 className={styles.comment}>{obj.comment}</h3>
      <div className={styles.detailscore}>
        <div className={`${styles.stageinfo} ${styles.stage}`}>
          {data.map((v, i) => {
            return (
              <div key={i} className={styles.onestage}>
                <div className={styles.stagename}>{v.pairname}</div>
                <div className={styles.stagedate}>{v.date}</div>
                <div
                  className={styles.stageroe}
                  style={
                    v.roe > 0 ? { color: "#26a69a" } : { color: "#ef5350" }
                  }
                >
                  {v.roe.toFixed(2)} %
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default MoreInfo;
