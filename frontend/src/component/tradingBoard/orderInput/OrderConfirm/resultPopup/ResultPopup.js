import React, { useState } from "react";
import styles from "./ResultPopup.module.css";
import VerticalLine from "../../../../lines/VerticalLine";
import HorizontalLine from "../../../../lines/HorizontalLine";
import { getAuth } from "firebase/auth";

const ResultPopup = (props) => {
  const auth = getAuth();

  const goRanking = () => {
    fetch("http://43.202.77.76:5000/ranking", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        // TODO: update userid for firebase
        user: auth.currentUser.uid,
        displayname: auth.currentUser.displayName,
        photourl: auth.currentUser.photoURL,
        scoreid: props.scoreid,
        balance: props.balance,
      }),
    }).then(window.location.replace("/ranking"));
  };
  const retry = () => {
    window.location.reload();
  };
  return (
    <div className={styles.modal}>
      {props.submitOrder ? (
        <div className={styles.result}>Fast Forwarding...</div>
      ) : (
        <div className={styles.result}>
          <div className={styles.header}>
            <div className={styles.headertitle}>
              <div className={styles.headerentry}>{props.result.entrytime}</div>
              <div className={styles.headername}>{props.result.name}</div>
            </div>
            <div
              className={styles.headerlev}
              style={{ color: `${props.color}` }}
            >
              X{props.result.leverage}
            </div>
          </div>
          <HorizontalLine />
          <div
            className={styles.roe}
            style={
              props.result.roe > 0 ? { color: "#26a69a" } : { color: "#ef5350" }
            }
          >
            {Math.floor(100 * (props.result.roe - props.leverage * 0.02)) / 100}{" "}
            %
          </div>
          <div className={styles.horizontalfield}>
            <div className={styles.infovalue} title={"PNL + Commisison"}>
              {Math.floor((props.result.pnl - props.result.commission) * 100) /
                100}{" "}
              USDT
            </div>
            <VerticalLine className={styles.vertical} />
            <div className={styles.infovalue}>+ {props.result.outtime} H</div>
          </div>
          {props.result.isliquidated ? (
            <div className={styles.liquidated}>포지션이 청산 되었습니다.</div>
          ) : null}

          <div className={styles.buttonfield}>
            {props.result.stage < 10 && !props.result.isliquidated ? (
              <button
                onClick={props.close}
                disabled={props.submitOrder ? true : false}
              >
                NEXT
              </button>
            ) : props.result.stage === 10 ? (
              props.mode === "competition" ? (
                <button
                  onClick={goRanking}
                  disabled={props.submitOrder ? true : false}
                >
                  스코어 등재하기
                </button>
              ) : (
                <button
                  onClick={retry}
                  disabled={props.submitOrder ? true : false}
                >
                  RETRY
                </button>
              )
            ) : (
              <button
                onClick={retry}
                disabled={props.submitOrder ? true : false}
              >
                RETRY
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
export default ResultPopup;
