import { useState } from "react";
import PostOrderJson from "../../../backendConn/PostOrderJson";
import ResultPopup from "./resultPopup/ResultPopup";
import styles from "./OrderConfirm.module.css";

function Orderconfirm({
  order,
  back,
  submitOrder,
  setSubmitOrder,
  orderInit,
  setPairtitle,
  setCandles,
  setResultChart,
  setResultScore,
  setIndex,
  balance,
  setBalance,
  setTitleaArray,
  color,
}) {
  const [receivedScore, setReceivedScore] = useState({
    stage: 0,
    name: "",
    entry_time: "",
    leverage: 0,
    entry_price: 0,
    end_price: 0,
    out_time: 0,
    roe: 0,
    pnl: 0,
    commission: 0,
    is_liquidated: false,
  });
  var profitROE;
  var lossROE;

  if (order.is_long) {
    profitROE = (order.profit_price - order.entry_price) / order.entry_price;
    lossROE = (order.entry_price - order.loss_price) / order.entry_price;
  } else {
    profitROE = (order.entry_price - order.profit_price) / order.entry_price;
    lossROE = (order.loss_price - order.entry_price) / order.entry_price;
  }

  const [modalOpen, setModalOpen] = useState(false);
  const [invalidOrder, setinvalidOrder] = useState(false);

  const profitPNL = order.entry_price * order.quantity * profitROE;
  const lossPNL = order.entry_price * order.quantity * lossROE;

  const backClick = () => {
    back((current) => !current);
  };
  const closeModal = () => {
    setModalOpen(false);
    back((current) => !current);
    setIndex((current) => current + 1);
  };
  const finalConfirm = () => {
    const resultPromise = PostOrderJson(
      "http://bitmoi.co.kr:5000/" + order.mode,
      order
    );
    resultPromise
      .then((rchart) => {
        if (order.mode === "competition") {
          setPairtitle(rchart.score.name);
          setTitleaArray((current) => [
            ...current,
            rchart.resultscore.name + ",",
          ]);
          setCandles(rchart.origin_chart);
        }
        setResultChart(rchart.result_chart);
        setResultScore(rchart.score);
        setBalance(
          (current) => current + rchart.score.pnl - rchart.score.commission
        );
        setReceivedScore(rchart.score);
      })
      .catch((error) => {
        console.log(error);
        setinvalidOrder(true);
        return;
      })
      .then(() => {
        setSubmitOrder(true);
        setModalOpen(true);
        setinvalidOrder(false);
      });
    setTimeout(() => {
      orderInit();
    }, 2390);
  };

  return (
    <div className={styles.confirmwindow}>
      {modalOpen ? (
        <ResultPopup
          close={closeModal}
          result={receivedScore}
          mode={order.mode}
          submitOrder={submitOrder}
          color={color}
          balance={balance}
          scoreid={order.score_id}
          leverage={order.leverage}
        />
      ) : (
        <div className={styles.orderconfirm}>
          <button onClick={backClick} className={styles.backbutton}>
            돌아가기
          </button>
          <div className={styles.confirmtitle}>주문 확인</div>
          <div className={styles.confirmbody}>
            <div>
              현 진입 시점으로부터 24시간 안에 시장 가격이{" "}
              <span className={styles.highlight}>
                {order.profit_price > 0
                  ? order.profit_price.toLocaleString("en-US", {
                      maximumFractionDigits: 4,
                    })
                  : ""}{" "}
                USDT
              </span>
              에 도달하면
              <span className={styles.profit}>
                {" "}
                {profitPNL.toLocaleString("en-US", {
                  maximumFractionDigits: 2,
                })}{" "}
                USDT{" "}
              </span>
              만큼 수익을 실현합니다.
            </div>
            <div>
              반대로{" "}
              <span className={styles.highlight}>
                {order.loss_price > 0
                  ? order.loss_price.toLocaleString("en-US", {
                      maximumFractionDigits: 4,
                    })
                  : ""}{" "}
                USDT
              </span>
              에 도달하면
              <span className={styles.loss}>
                {" "}
                -
                {lossPNL.toLocaleString("en-US", {
                  maximumFractionDigits: 2,
                })}{" "}
                USDT{" "}
              </span>
              손절매를 합니다.
            </div>
            <div>
              만약 이 가격들에 도달하지 못하여 예약 주문이 체결되지 않을 경우,{" "}
              <span className={styles.highlight}>24시간 뒤 포지션을 정리</span>
              하고 가격 차이만큼 수익 또는 손실을 실현합니다.
            </div>
          </div>

          <div className={styles.submitbutton}>
            <button
              onClick={finalConfirm}
              className={
                order.isLong
                  ? `${styles.confirmlong}`
                  : `${styles.confirmshort}`
              }
            >
              주문 제출하기
            </button>
          </div>
          {invalidOrder ? (
            <div className={styles.invalidorder}>잘못된 주문입니다.</div>
          ) : (
            <div></div>
          )}
        </div>
      )}
    </div>
  );
}

export default Orderconfirm;
