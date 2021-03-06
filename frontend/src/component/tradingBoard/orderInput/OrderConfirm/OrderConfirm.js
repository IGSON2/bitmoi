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
    entrytime: "",
    name: "",
    leverage: 0,
    entryprice: 0,
    outhour: 0,
    roe: 0,
    pnl: 0,
    commission: 0,
    isliquidated: false,
  });
  var profitROE;
  var lossROE;
  if (order.isLong) {
    profitROE = (order.profitPrice - order.entryPrice) / order.entryPrice;
    lossROE = (order.entryPrice - order.lossPrice) / order.entryPrice;
  } else {
    profitROE = (order.entryPrice - order.profitPrice) / order.entryPrice;
    lossROE = (order.lossPrice - order.entryPrice) / order.entryPrice;
  }

  const [modalOpen, setModalOpen] = useState(false);

  const profitPNL = order.entryPrice * order.quantity * profitROE;
  const lossPNL = order.entryPrice * order.quantity * lossROE;

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
      "http://www.bitmoi.net/api/" + order.mode,
      order
    );
    resultPromise
      .then((rchart) => {
        if (order.mode === "competition") {
          setPairtitle(rchart.resultscore.name);
          setTitleaArray((current) => [
            ...current,
            rchart.resultscore.name + ",",
          ]);
          setCandles(rchart.originchart);
        }
        setResultChart(rchart.resultchart);
        setResultScore(rchart.resultscore);
        setBalance(
          (current) =>
            current + rchart.resultscore.pnl - rchart.resultscore.commission
        );
        setReceivedScore(rchart.resultscore);
      })
      .then(setSubmitOrder(true))
      .then(setModalOpen(true));
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
          scoreid={order.scoreid}
          leverage={order.leverage}
        />
      ) : (
        <div className={styles.orderconfirm}>
          <button onClick={backClick} className={styles.backbutton}>
            ????????????
          </button>
          <div className={styles.confirmtitle}>?????? ??????</div>
          <div className={styles.confirmbody}>
            <div>
              ??? ?????? ?????????????????? 24?????? ?????? ?????? ?????????{" "}
              <span className={styles.highlight}>
                {order.profitPrice > 0
                  ? order.profitPrice.toLocaleString("en-US", {
                      maximumFractionDigits: 4,
                    })
                  : ""}{" "}
                USDT
              </span>
              ??? ????????????
              <span className={styles.profit}>
                {" "}
                {profitPNL.toLocaleString("en-US", {
                  maximumFractionDigits: 2,
                })}{" "}
                USDT{" "}
              </span>
              ?????? ????????? ???????????????.
            </div>
            <div>
              ?????????{" "}
              <span className={styles.highlight}>
                {order.lossPrice > 0
                  ? order.lossPrice.toLocaleString("en-US", {
                      maximumFractionDigits: 4,
                    })
                  : ""}{" "}
                USDT
              </span>
              ??? ????????????
              <span className={styles.loss}>
                {" "}
                -
                {lossPNL.toLocaleString("en-US", {
                  maximumFractionDigits: 2,
                })}{" "}
                USDT{" "}
              </span>
              ???????????? ?????????.
            </div>
            <div>
              ?????? ??? ???????????? ???????????? ????????? ?????? ????????? ???????????? ?????? ??????,{" "}
              <span className={styles.highlight}>24?????? ??? ???????????? ??????</span>
              ?????? ?????? ???????????? ?????? ?????? ????????? ???????????????.
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
              ?????? ????????????
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default Orderconfirm;
