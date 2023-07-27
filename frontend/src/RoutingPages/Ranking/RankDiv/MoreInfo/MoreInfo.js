import styles from "./MoreInfo.module.css";
import { useState } from "react";
import Timeformatter from "../../../../component/Timeformatter/Timeformatter";
import WordConfirm from "./WordConfirm/WordConfirm";

function MoreInfo({ setMoreInfo, data, certified }) {
  const [myword, setMyword] = useState("");
  const [popupOpen, setPopupOpen] = useState(false);
  const mywordChange = (e) => {
    setMyword(e.target.value);
  };
  const wordSubmit = (e) => {
    e.preventDefault();
    setPopupOpen(true);
  };

  return (
    <div className={styles.moreinfo}>
      {popupOpen ? (
        <WordConfirm popupOpen={setPopupOpen} comment={myword} />
      ) : null}
      <h3
        className={styles.comment}
        style={
          certified
            ? { backgroundColor: "#faebef" }
            : data.comment === ""
            ? { backgroundColor: "transparent" }
            : null
        }
      >
        {data.comment === "" ? (
          certified ? (
            <form className={styles.myword} onSubmit={wordSubmit}>
              <input
                type="text"
                placeholder="기록 갱신을 축하 드립니다.  소감을 등록해 주세요."
                value={myword}
                maxLength="80"
                onChange={mywordChange}
              ></input>
              <button>등록하기</button>
            </form>
          ) : (
            `아직 등록된 소감이 없습니다.`
          )
        ) : (
          `" ${data.comment} "`
        )}
      </h3>
      <div className={styles.detailscore}>
        <div className={`${styles.field} ${styles.date}`}>
          <div className={styles.title}>등재 일자</div>
          <div className={styles.value} style={{ letterSpacing: "1px" }}>
            {Timeformatter(data.scoreid)}
          </div>
        </div>
        <div className={`${styles.field} ${styles.lev}`}>
          <div className={styles.title}>평균 레버리지</div>
          <div className={styles.value}>X {data.avglev}</div>
        </div>
        <div className={`${styles.field} ${styles.pnl}`}>
          <div className={styles.title}>평균 일당</div>
          <div
            className={styles.value}
            style={
              data.avgpnl > 0 ? { color: "#26a69a" } : { color: "#ef5350" }
            }
          >
            $ {data.avgpnl}
          </div>
        </div>
        <div className={`${styles.field} ${styles.roe}`}>
          <div className={styles.title}>평균 수익률</div>
          <div
            className={styles.value}
            style={
              data.avgroe > 0 ? { color: "#26a69a" } : { color: "#ef5350" }
            }
          >
            {data.avgroe} %
          </div>
        </div>
        <div className={`${styles.stageinfo} ${styles.stage}`}>
          {data.stagearray.map((v, i) => {
            return (
              <div
                key={i}
                className={styles.onestage}
                style={
                  v.name.includes("(Liq.)")
                    ? { background: "black", color: "white" }
                    : null
                }
              >
                <div className={styles.stagename}>
                  {v.name.includes("(Liq.)")
                    ? v.name.replace(" (Liq.)", "")
                    : v.name}
                </div>
                <div className={styles.stagedate}>{v.date}</div>
                <div
                  className={styles.stageroe}
                  style={
                    v.roe > 0 ? { color: "#26a69a" } : { color: "#ef5350" }
                  }
                >
                  {v.roe} %
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
