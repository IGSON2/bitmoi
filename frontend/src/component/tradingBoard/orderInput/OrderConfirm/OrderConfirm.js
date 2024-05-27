import { useState } from "react";
import ResultPopup from "./resultPopup/ResultPopup";
import styles from "./OrderConfirm.module.css";
import axiosClient from "../../../backendConn/axiosClient";

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
  setTitleArray,
  color,
  userInfo,
}) {
  const [isChecked, setIsChecked] = useState(true);
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

  const finalConfirm = async () => {
    try {
      const response = await axiosClient.post(`/basic/${order.mode}`, order);
      if (order.mode === "competition") {
        setPairtitle(response.data.score.name);
        setTitleArray((current) => [...current, response.data.score.name]);
        response.data.origin_chart.pdata.reverse();
        response.data.origin_chart.vdata.reverse();
        setCandles(response.data.origin_chart);
      }
      setResultChart(response.data.result_chart);
      setResultScore(response.data.score);
      setBalance(
        (current) =>
          current + response.data.score.pnl - response.data.score.commission
      );
      setReceivedScore(response.data.score);
      setSubmitOrder(true);
      setModalOpen(true);
      setinvalidOrder(false);
    } catch (error) {
      console.error(error);
      setinvalidOrder(true);
      return;
    }

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
          order={order}
          submitOrder={submitOrder}
          color={color}
          balance={balance}
          userInfo={userInfo}
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
            {order.mode === "competition" && order.stage === 1 ? (
              <label className={styles.checkbox}>
                <input
                  type="checkbox"
                  checked={isChecked}
                  style={{ width: "15px", height: "15px" }}
                  onChange={(e) => setIsChecked(e.target.checked)}
                />{" "}
                주문을 제출하면 MOI 토큰이 1개 차감되고 도전이 시작됩니다.
                <br />
              </label>
            ) : null}

            <button
              onClick={finalConfirm}
              disabled={order.stage === 1 ? !isChecked : false}
              className={
                order.is_long
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
