import React, { useState } from "react";
import styles from "./ResultPopup.module.css";
import VerticalLine from "../../../../lines/VerticalLine";
import HorizontalLine from "../../../../lines/HorizontalLine";
import axiosClient from "../../../../backendConn/axiosClient";
import { BsXLg } from "react-icons/bs";

const ResultPopup = (props) => {
  const [comment, setComment] = useState("");
  const [openComment, setOpenComment] = useState(false);

  const goRanking = async () => {
    try {
      const response = await axiosClient.post("/rank", {
        comment: comment,
        score_id: props.order.score_id,
      });
      if (response.status === 200) {
        window.location.replace("/rank/1");
      } else {
        throw new Error(response.data);
      }
    } catch (error) {
      if (error.response.data.includes("low score")) {
        // 기록 갱신 실패
      }
      console.error(error);
    }
  };

  const retry = () => {
    window.location.reload();
  };

  const closePopup = () => {
    setOpenComment(false);
  };

  const openPopup = () => {
    setOpenComment(true);
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
              X{props.order.leverage}
            </div>
          </div>
          <HorizontalLine />
          <div
            className={styles.roe}
            style={
              props.result.roe > 0 ? { color: "#26a69a" } : { color: "#ef5350" }
            }
          >
            {Math.floor(
              100 * (props.result.roe - props.order.leverage * 0.02)
            ) / 100}{" "}
            %
          </div>
          <div className={styles.horizontalfield}>
            <div className={styles.infovalue} title={"PNL + Commisison"}>
              {Math.floor((props.result.pnl - props.result.commission) * 100) /
                100}{" "}
              USDT
            </div>
            <VerticalLine className={styles.vertical} />
            <div className={styles.infovalue}>+ {props.result.out_time} H</div>
          </div>
          {props.result.is_liquidated ? (
            <div className={styles.liquidated}>포지션이 청산 되었습니다.</div>
          ) : null}

          <div className={styles.buttonfield}>
            {props.result.stage < 10 && !props.result.is_liquidated ? (
              <button
                onClick={props.close}
                disabled={props.submitOrder ? true : false}
              >
                NEXT
              </button>
            ) : props.result.stage === 10 ? (
              props.order.mode === "competition" ? (
                <button
                  onClick={openPopup}
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
      {openComment ? (
        <div className={styles.comment}>
          <div className={styles.background} onClick={closePopup}></div>
          <div className={styles.inner}>
            <div className={styles.closebutton}>
              <span>
                <BsXLg onClick={closePopup} />
              </span>
            </div>
            <div className={styles.title}>당신의 총 스코어</div>
            <div className={styles.score}>
              {props.balance.toFixed(2)}
              <span>USDT</span>
            </div>
            <input
              className={styles.textinput}
              type="text"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="소감을 입력해주세요."
              maxLength={100}
            />
            <button className={styles.sendbutton} onClick={goRanking}>
              등록하기
            </button>
          </div>
        </div>
      ) : null}
    </div>
  );
};
export default ResultPopup;
