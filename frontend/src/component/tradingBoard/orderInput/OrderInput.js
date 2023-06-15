import React, { useEffect, useRef, useState } from "react";
import styles from "./OrderInput.module.css";
import Warning from "./Warning";
import OrderConfirm from "./OrderConfirm/OrderConfirm";
import { getAuth, onAuthStateChanged } from "firebase/auth";
import { AiOutlineCloseCircle, AiOutlinePlusCircle } from "react-icons/ai";

function OrderInput({
  mode,
  name,
  index,
  opened,
  submitOrder,
  setSubmitOrder,
  entryPrice,
  setProfitMarker,
  setLossMarker,
  identifier,
  setName,
  setCandles,
  setResultChart,
  setIndex,
  setResultScore,
  balance,
  setBalance,
  setTitleaArray,
  entryTime,
  resultTerm,
}) {
  const auth = getAuth();
  const [quantity, setQuantity] = useState();
  const [quantityRate, setQuantityRate] = useState(1);
  const [profitPrice, setProfitPrice] = useState();
  const [profitRate, setProfitRate] = useState(1);
  const [lossPrice, setLossPrice] = useState();
  const [lossRate, setLossRate] = useState(1);
  const [leverage, setLeverage] = useState(1);
  const startBalance = 1000;
  const [isLong, setIsLong] = useState(true);
  const [color, setColor] = useState("black");
  const [levInputMode, setLevInputMode] = useState(false);
  const [profitLossErr, setProfitLossErr] = useState(true);
  const [quanErr, setQuanErr] = useState(true);
  const [compLoginErr, setCompLoginErr] = useState(true);
  const [loginWarning, setLoginWarning] = useState("");
  const [quanWarning, setQuanWarning] = useState("");
  const [profitWarning, setProfitWarning] = useState("");
  const [lossWarning, setLossWarning] = useState("");
  const [levWarning, setLevWarning] = useState("");
  const [confirm, setConfirm] = useState(false);
  const [quanClosed, setQuanClosed] = useState(false);
  const [profitClosed, setProfitClosed] = useState(false);
  const [lossClosed, setLossClosed] = useState(false);
  const [orderObject, setOrderObject] = useState({
    mode: "",
    user_id: "",
    name: "",
    entrytime: "",
    stage: 0,
    isLong: true,
    entry_price: 0,
    quantity: 0,
    quantity_rate: 0,
    profit_price: 0,
    loss_price: 0,
    leverage: 0,
    balance: 0,
    identifier: "",
    score_id: "",
    resultterm: 0,
  });
  var commission = 0.0002;
  const inputRef = useRef(null);
  const longClicked = () => {
    setIsLong(true);
    setProfitWarning("");
    setLossWarning("");
    setLevWarning("");
    setQuanWarning("");
  };
  const shortClicked = () => {
    setIsLong(false);
    setProfitWarning("");
    setLossWarning("");
    setLevWarning("");
    setQuanWarning("");
  };
  const onsubmit = (event) => {
    event.preventDefault();
    var tempObject = {
      mode: mode,
      user_id: "",
      name: name,
      entrytime: entryTime,
      stage: index + 1,
      isLong: isLong,
      entry_price: entryPrice,
      quantity: quantity,
      quantity_rate: quantityRate,
      profit_price: profitPrice,
      loss_price: lossPrice,
      leverage: leverage,
      balance: balance,
      identifier: identifier,
      score_id: Date.now().toString(),
      resultTerm: resultTerm,
    };
    // Update userid for firebase
    if (mode === "competition") {
      tempObject.user_id = auth.currentUser.uid;
    }
    setOrderObject(tempObject);
    setConfirm((current) => !current);
    setLossMarker(lossPrice);
    setProfitMarker(profitPrice);
  };
  const orderInit = () => {
    setIsLong(true);
    setQuantity();
    setProfitPrice();
    setLossPrice();
    setLeverage(1);
    setQuantityRate(1);
    setProfitRate(1);
    setLossRate(1);
    setProfitWarning("");
    setLossWarning("");
    setLevWarning("");
    setQuanWarning("");
  };
  const enterBlock = (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
    }
  };

  const quantityChange = (event) => {
    if (event.target.valueAsNumber < 0) {
      setQuantity(0);
    } else {
      setQuantity(event.target.valueAsNumber);
    }
  };
  const quantityRateChange = (event) => {
    setQuantityRate(event.target.valueAsNumber);
    setQuantity(
      Math.floor(
        ((balance * leverage * 0.9998) / entryPrice) *
          (event.target.valueAsNumber / 100) *
          10000
      ) / 10000
    );
  };
  const profitChange = (event) => {
    if (event.target.valueAsNumber < 0) {
      setProfitPrice(0);
    } else {
      setProfitPrice(event.target.valueAsNumber);
      if (isLong) {
        setProfitRate(
          Math.floor(
            (10000 * (event.target.valueAsNumber - entryPrice)) / entryPrice
          ) / 100
        );
      } else {
        setProfitRate(
          Math.floor(
            (10000 * (entryPrice - event.target.valueAsNumber)) / entryPrice
          ) / 100
        );
      }
    }
  };
  const profitRateChange = (event) => {
    setProfitRate(event.target.valueAsNumber);
    if (isLong) {
      setProfitPrice(
        Math.floor(
          entryPrice * (1 + event.target.valueAsNumber / 100) * 10000
        ) / 10000
      );
    } else {
      setProfitPrice(
        Math.floor(
          entryPrice * (1 - event.target.valueAsNumber / 100) * 10000
        ) / 10000
      );
      if (event.target.valueAsNumber >= 100) {
        setProfitPrice(Math.floor(entryPrice * (1 - 0.9999) * 10000) / 10000);
      }
    }
  };
  const lossChange = (event) => {
    if (event.target.valueAsNumber < 0) {
      setLossPrice(0);
    } else {
      setLossPrice(event.target.valueAsNumber);
      if (isLong) {
        setLossRate(
          Math.floor(
            (10000 * (entryPrice - event.target.valueAsNumber)) / entryPrice
          ) / 100
        );
      } else {
        setLossRate(
          Math.floor(
            (10000 * (event.target.valueAsNumber - entryPrice)) / entryPrice
          ) / 100
        );
      }
    }
  };
  const lossRateChange = (event) => {
    setLossRate(event.target.valueAsNumber);
    if (isLong) {
      setLossPrice(
        Math.ceil(entryPrice * (1 - event.target.valueAsNumber / 100) * 10000) /
          10000
      );
    } else {
      setLossPrice(
        Math.floor(
          entryPrice * (1 + event.target.valueAsNumber / 100) * 10000
        ) / 10000
      );
    }
  };
  const leverageChange = (event) => {
    if (!event.target.valueAsNumber) {
      setLeverage("");
    } else {
      setLeverage(event.target.valueAsNumber);
      setQuantity(
        Math.floor(
          ((balance * event.target.valueAsNumber * 0.9998) / entryPrice) *
            (quantityRate / 100) *
            10000
        ) / 10000
      );
    }
    if (event.target.valueAsNumber === 1) {
      setColor("black");
    } else {
      setColor(
        `rgb(
          ${Math.round(2.5 * event.target.valueAsNumber) - 1},
          ${30 + Math.round(0.53 * event.target.valueAsNumber) - 1},
        ${167 - Math.round(0.71 * event.target.valueAsNumber) - 1}
        )`
      );
    }
  };
  const levInputChange = () => {
    setLevInputMode((current) => !current);
  };
  const onWheelQuan = (event) => {
    if (Number(event.deltaY) < 0) {
      if (quantityRate < 95) {
        setQuantityRate((current) => current + 5);
        setQuantity(
          Math.floor(
            ((balance * leverage * 0.9998) / entryPrice) *
              ((quantityRate + 5) / 100) *
              10000
          ) / 10000
        );
      } else {
        setQuantityRate(100);
        setQuantity(
          Math.floor(((balance * leverage * 0.9998) / entryPrice) * 10000) /
            10000
        );
      }
    } else {
      if (quantityRate > 0) {
        setQuantityRate((current) => current - 5);
        setQuantity(
          Math.floor(
            ((balance * leverage * 0.9998) / entryPrice) *
              ((quantityRate - 5) / 100) *
              10000
          ) / 10000
        );
      } else {
        setQuantityRate(0);
        setQuantity(0);
      }
    }
  };
  const onWheelProfit = (event) => {
    if (Number(event.deltaY) < 0) {
      if (profitRate < 99.5) {
        setProfitRate((current) => current + 0.5);
        if (isLong) {
          setProfitPrice(
            Math.floor(entryPrice * (1 + (profitRate + 0.5) / 100) * 10000) /
              10000
          );
        } else {
          setProfitPrice(
            Math.floor(entryPrice * (1 - (profitRate + 0.5) / 100) * 10000) /
              10000
          );
        }
      } else {
        setProfitRate(99.5);
        if (isLong) {
          setProfitPrice(Math.floor(entryPrice * 2 * 10000) / 10000);
        } else {
          setProfitPrice(
            Math.floor(entryPrice * (1 - 99.5 / 100) * 10000) / 10000
          );
        }
      }
    } else {
      if (profitRate > 0) {
        setProfitRate((current) => current - 0.5);
        if (isLong) {
          setProfitPrice(
            Math.floor(entryPrice * (1 + (profitRate - 0.5) / 100) * 10000) /
              10000
          );
        } else {
          setProfitPrice(
            Math.floor(entryPrice * (1 - (profitRate - 0.5) / 100) * 10000) /
              10000
          );
        }
      } else {
        setProfitRate(0);
        setProfitPrice(Math.floor(entryPrice * 10000) / 10000);
      }
    }
  };
  const onWheelLoss = (event) => {
    if (Number(event.deltaY) < 0) {
      if (lossRate < 99.5) {
        setLossRate((current) => current + 0.5);
        if (isLong) {
          setLossPrice(
            Math.ceil(entryPrice * (1 - (lossRate + 0.5) / 100) * 10000) / 10000
          );
        } else {
          setLossPrice(
            Math.floor(entryPrice * (1 + (lossRate + 0.5) / 100) * 10000) /
              10000
          );
        }
      } else {
        setLossRate(99.5);
        if (isLong) {
          setLossPrice(Math.ceil(entryPrice * (0.5 / 100) * 10000) / 10000);
        } else {
          setLossPrice(
            Math.floor(entryPrice * (1 + 99.5 / 100) * 10000) / 10000
          );
        }
      }
    } else {
      if (lossRate > 0) {
        setLossRate((current) => current - 0.5);
        if (isLong) {
          setLossPrice(
            Math.ceil(entryPrice * (1 - (lossRate - 0.5) / 100) * 10000) / 10000
          );
        } else {
          setLossPrice(
            Math.floor(entryPrice * (1 + (lossRate - 0.5) / 100) * 10000) /
              10000
          );
        }
      } else {
        setLossRate(0);
        setLossPrice(Math.ceil(entryPrice * 10000) / 10000);
      }
    }
  };
  const onWheelLev = (event) => {
    if (Number(event.deltaY) < 0) {
      if (leverage < 99) {
        setLeverage((current) => current + 1);
        setQuantity(
          Math.floor(
            ((balance * (leverage + 1) * 0.9998) / entryPrice) *
              (quantityRate / 100) *
              10000
          ) / 10000
        );
      } else {
        setLeverage(100);
        setQuantity(
          Math.floor(
            ((balance * 100 * 0.9998) / entryPrice) *
              (quantityRate / 100) *
              10000
          ) / 10000
        );
      }
    } else {
      if (leverage > 1) {
        setLeverage((current) => current - 1);
        setQuantity(
          Math.floor(
            ((balance * (leverage - 1) * 0.9998) / entryPrice) *
              (quantityRate / 100) *
              10000
          ) / 10000
        );
      } else {
        setLeverage(1);
        setQuantity(
          Math.floor(
            ((balance * 1 * 0.9998) / entryPrice) * (quantityRate / 100) * 10000
          ) / 10000
        );
      }
    }

    setColor(
      `rgb(
          ${Math.round(2.5 * leverage) - 1},
          ${30 + Math.round(0.53 * leverage) - 1},
        ${167 - Math.round(0.71 * leverage) - 1}
        )`
    );
  };

  useEffect(() => {
    if (levInputMode) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [levInputMode]);

  React.useEffect(() => {
    if (isLong) {
      if (profitClosed) {
        setProfitPrice(Math.floor(entryPrice * 11 * 10000) / 10000);
      }
      if (lossClosed) {
        setLossPrice(
          Math.ceil((entryPrice - entryPrice * (0.98 / leverage)) * 100000) /
            100000
        );
      }
      if (entryPrice < profitPrice && entryPrice > lossPrice) {
        setProfitWarning("");
        setLossWarning("");
        setProfitLossErr(false);
        if (lossPrice) {
          if (
            0 < (entryPrice - lossPrice) / entryPrice &&
            (entryPrice - lossPrice) / entryPrice <= 1.00002 / leverage
          ) {
            setProfitLossErr(false);
            setLevWarning("");
          } else {
            setProfitLossErr(true);
            setLevWarning(
              `${leverage}배 레버리지의 최대 손실 감수 가격은` +
                "\n" +
                ` ${(entryPrice - entryPrice * (1 / leverage)).toFixed(
                  5
                )} USDT( -${(100 * (1 / leverage)).toFixed(2)}% )입니다.`
            );
          }
        } else {
          setProfitLossErr(true);
        }
      } else {
        setProfitLossErr(true);
        if (profitPrice && profitPrice <= entryPrice) {
          setProfitWarning("수익 실현 가격을 진입 가격보다 높게 설정해 주세요");
        } else {
          setProfitWarning("");
        }
        if (lossPrice && lossPrice >= entryPrice) {
          setLossWarning("손실 감수 가격을 진입 가격보다 낮게 설정해 주세요.");
        } else {
          setLossWarning("");
        }
      }
    } else {
      if (profitClosed) {
        setProfitPrice(Math.floor(entryPrice * 0.00001 * 10000) / 10000);
      }
      if (lossClosed) {
        setLossPrice(
          Math.floor((entryPrice + entryPrice * (0.98 / leverage)) * 10000) /
            10000
        );
      }
      if (entryPrice > profitPrice && entryPrice < lossPrice) {
        setProfitLossErr(false);
        setProfitWarning("");
        setLossWarning("");
        if (lossPrice) {
          if (
            0 < (lossPrice - entryPrice) / entryPrice &&
            (lossPrice - entryPrice) / entryPrice <= 1 / leverage
          ) {
            setProfitLossErr(false);
            setLevWarning("");
          } else {
            setProfitLossErr(true);
            setLevWarning(
              `${leverage}배 레버리지의 최대 손실 감수 가격은` +
                "\n" +
                `${(entryPrice + entryPrice * (1 / leverage)).toFixed(
                  5
                )} USDT( +${(100 * (1 / leverage)).toFixed(2)}% )입니다.`
            );
          }
        } else {
          setProfitLossErr(true);
        }
      } else {
        setProfitLossErr(true);
        if (profitPrice && profitPrice >= entryPrice) {
          setProfitWarning("수익 실현 가격을 진입 가격보다 낮게 설정해 주세요");
        } else {
          setProfitWarning("");
        }
        if (lossPrice && lossPrice <= entryPrice) {
          setLossWarning("손실 감수 가격을 진입 가격보다 높게 설정해 주세요.");
        } else {
          setLossWarning("");
        }
      }
    }
  }, [profitPrice, lossPrice, leverage, isLong]);

  useEffect(() => {
    if (quantity) {
      if (balance * leverage > entryPrice * quantity * (1 + commission)) {
        setQuanErr(false);
        setQuanWarning("");
      } else {
        setQuanErr(true);
        setQuanWarning(
          `잔고가 부족합니다.` +
            "\n" +
            `주문 가능한 최대 수량은 ${(
              (balance * leverage) /
              (entryPrice * (1 + commission))
            ).toFixed(5)} EA 입니다.`
        );
      }
    } else {
      setQuanErr(true);
    }
  }, [quantity, leverage, profitPrice, lossPrice]);

  useEffect(() => {
    if (mode == "competition") {
      onAuthStateChanged(auth, (user) => {
        if (user) {
          setCompLoginErr(false);
          setLoginWarning("");
        } else {
          setCompLoginErr(true);
          setLoginWarning("경쟁모드는 로그인이 필요한 서비스입니다.");
        }
      });
    } else if (mode === "practice") {
      setCompLoginErr(false);
      setLoginWarning("");
    }
  }, []);

  const quanClose = () => {
    setQuantityRate(100);
    setQuantity(
      Math.floor(((balance * leverage * 0.9998) / entryPrice) * 10000) / 10000
    );
    setQuanClosed(true);
  };

  const profitClose = () => {
    if (isLong) {
      setProfitPrice(Math.floor(entryPrice * 11 * 10000) / 10000);
    } else {
      setProfitPrice(Math.floor(entryPrice * 0.00001 * 10000) / 10000);
    }
    setProfitClosed(true);
  };

  const lossClose = () => {
    if (isLong) {
      setLossPrice(
        Math.floor((entryPrice - entryPrice * (0.99 / leverage)) * 10000) /
          10000
      );
    } else {
      setLossPrice(
        Math.floor((entryPrice + entryPrice * (0.99 / leverage)) * 10000) /
          10000
      );
    }
    setLossClosed(true);
  };

  const quanOpen = () => {
    setQuanClosed(false);
    setQuantity(
      Math.floor(((balance * leverage * 0.1) / entryPrice) * 10000) / 10000
    );
    setQuantityRate(10);
  };

  const profitOpen = () => {
    setProfitClosed(false);
    setProfitRate(10);
    if (isLong) {
      setProfitPrice(Math.floor(entryPrice * 1.1 * 10000) / 10000);
    } else {
      setProfitPrice(Math.floor(entryPrice * 0.9 * 10000) / 10000);
    }
  };

  const lossOpen = () => {
    setLossClosed(false);
    setLossRate(10);
    if (isLong) {
      setLossPrice(Math.floor((entryPrice - entryPrice * 0.1) * 10000) / 10000);
    } else {
      setLossPrice(Math.floor((entryPrice + entryPrice * 0.1) * 10000) / 10000);
    }
  };

  return (
    <div
      className={`${opened ? styles.ordernavopened : styles.ordernavclosed}`}
    >
      {confirm ? (
        <div className={styles.submitornot}>
          <OrderConfirm
            order={orderObject}
            back={setConfirm}
            submitOrder={submitOrder}
            setSubmitOrder={setSubmitOrder}
            orderInit={orderInit}
            setPairtitle={setName}
            setCandles={setCandles}
            setResultChart={setResultChart}
            setResultScore={setResultScore}
            setIndex={setIndex}
            balance={balance}
            setBalance={setBalance}
            setTitleaArray={setTitleaArray}
            color={color}
          />
        </div>
      ) : (
        <div className={styles.submitornot}>
          <div className={styles.positionselector}>
            <div className={styles.long} title={"가격 상승을 예측합니다."}>
              <button
                className={
                  isLong ? styles.positionbutton_active : styles.positionbutton
                }
                onClick={longClicked}
              >
                Long ↑
              </button>
            </div>
            <div className={styles.short} title={"가격 하락을 예측합니다."}>
              <button
                className={
                  isLong ? styles.positionbutton : styles.positionbutton_active
                }
                onClick={shortClicked}
              >
                Short ↓
              </button>
            </div>
          </div>
          <div className={styles.orderbox}>
            <div className={styles.orderheader}>
              <div className={styles.ordertitle}>ORDER</div>
              <div className={styles.stage}>
                {index < 9 ? "STAGE 0" + `${index + 1}` : "Last Stage"}
              </div>
            </div>
            <form
              className={styles.orderform}
              onKeyDown={enterBlock}
              onSubmit={onsubmit}
            >
              <div className={styles.lableinputs}>
                <label className={styles.lables}>Entry price</label>
                <div className={styles.entrypricediv}>
                  {entryPrice.toFixed(4)}
                </div>
                <div className={styles.inputrightdiv}>USDT</div>
              </div>

              {quanClosed ? (
                <div className={styles.lableinputs}>
                  <div className={styles.lables}>Quantity</div>
                  <div className={styles.maxvalue}>{quantity}</div>
                  <div className={styles.inputrightdiv}>MAX</div>
                  <div className={styles.inputclose} onClick={quanOpen}>
                    <AiOutlinePlusCircle />
                  </div>
                </div>
              ) : (
                <div className={styles.lableinputs} onWheel={onWheelQuan}>
                  <label className={styles.lables} htmlFor="quantity">
                    Quantity
                  </label>
                  <div className={styles.twoinputs}>
                    <input
                      className={quanWarning === "" ? "" : styles.quanwarning}
                      id="quantity"
                      type={"number"}
                      step={"0.0001"}
                      value={quantity}
                      onChange={quantityChange}
                      placeholder="수량을 입력하세요."
                    ></input>
                    <input
                      className={isLong ? styles.longinput : styles.shortinput}
                      min={0}
                      max={100}
                      type={"range"}
                      step={"5"}
                      value={quantityRate}
                      onChange={quantityRateChange}
                    ></input>
                  </div>

                  <label className={styles.inputrightdiv} htmlFor="quantity">
                    EA
                  </label>
                  <div className={styles.inputclose} onClick={quanClose}>
                    <AiOutlineCloseCircle />
                  </div>
                </div>
              )}

              {profitClosed ? (
                <div className={styles.lableinputs}>
                  <div className={styles.lables}>Take profit</div>
                  <div className={styles.maxvalue}>{profitPrice}</div>
                  <div className={styles.inputrightdiv}>MAX</div>
                  <div className={styles.inputclose} onClick={profitOpen}>
                    <AiOutlinePlusCircle />
                  </div>
                </div>
              ) : (
                <div className={styles.lableinputs} onWheel={onWheelProfit}>
                  <label className={styles.lables} htmlFor="profitprice">
                    Take profit
                  </label>
                  <div className={styles.twoinputs}>
                    <input
                      className={`${
                        profitWarning === "" ? "" : styles.profitwarning
                      }`}
                      id="profitprice"
                      type={"number"}
                      step={"0.0001"}
                      value={profitPrice}
                      onChange={profitChange}
                      placeholder="수익 실현 가격을 입력하세요."
                    ></input>
                    <input
                      className={isLong ? styles.longinput : styles.shortinput}
                      type={"range"}
                      min={0}
                      max={100}
                      step={"0.5"}
                      value={profitRate}
                      onChange={profitRateChange}
                    ></input>
                  </div>

                  <label htmlFor="profitprice" className={styles.inputrightdiv}>
                    <div className={styles.profitbox}>
                      {entryPrice > 0 && profitPrice > 0
                        ? isLong
                          ? (
                              ((profitPrice - entryPrice) / entryPrice) *
                              100
                            ).toFixed(2) + "%"
                          : (
                              ((entryPrice - profitPrice) / entryPrice) *
                              100
                            ).toFixed(2) + "%"
                        : "0%"}
                    </div>
                  </label>
                  <div className={styles.inputclose} onClick={profitClose}>
                    <AiOutlineCloseCircle />
                  </div>
                </div>
              )}

              {lossClosed ? (
                <div className={styles.lableinputs}>
                  <div className={styles.lables}>Stop loss</div>
                  <div className={styles.maxvalue}>{lossPrice}</div>
                  <div className={styles.inputrightdiv}>MAX</div>
                  <div className={styles.inputclose} onClick={lossOpen}>
                    <AiOutlinePlusCircle />
                  </div>
                </div>
              ) : (
                <div className={styles.lableinputs} onWheel={onWheelLoss}>
                  <label className={styles.lables} htmlFor="stoplossprice">
                    Stop loss
                  </label>
                  <div className={styles.twoinputs}>
                    <input
                      className={`${
                        lossWarning === ""
                          ? levWarning
                            ? styles.losswarning
                            : null
                          : styles.losswarning
                      }`}
                      id="stoplossprice"
                      type={"number"}
                      step={"0.0001"}
                      value={lossPrice}
                      onChange={lossChange}
                      placeholder="손실 감수 가격을 입력하세요."
                    ></input>
                    <input
                      className={isLong ? styles.longinput : styles.shortinput}
                      type={"range"}
                      min={0}
                      max={99.5}
                      step={"0.5"}
                      value={lossRate}
                      onChange={lossRateChange}
                    ></input>
                  </div>

                  <label
                    htmlFor="stoplossprice"
                    className={styles.inputrightdiv}
                  >
                    <div className={styles.lossbox}>
                      {entryPrice > 0 && lossPrice > 0
                        ? isLong
                          ? (
                              ((lossPrice - entryPrice) / entryPrice) *
                              100
                            ).toFixed(2) + "%"
                          : (
                              ((entryPrice - lossPrice) / entryPrice) *
                              100
                            ).toFixed(2) + "%"
                        : "0%"}
                    </div>
                  </label>
                  <div className={styles.inputclose} onClick={lossClose}>
                    <AiOutlineCloseCircle />
                  </div>
                </div>
              )}

              <div className={styles.levlableinput} onWheel={onWheelLev}>
                <label className={styles.lables} htmlFor="leverage">
                  Leverage
                </label>
                {levInputMode ? (
                  <input
                    className={styles.levinput}
                    id="leverage"
                    type={"number"}
                    value={leverage}
                    min={1}
                    max={100}
                    step={"1"}
                    onChange={leverageChange}
                    ref={inputRef}
                  ></input>
                ) : (
                  <input
                    className={`${styles.levrangeinput} ${
                      isLong ? styles.longinput : styles.shortinput
                    }`}
                    type={"range"}
                    value={leverage}
                    min={1}
                    max={100}
                    step={"1"}
                    onChange={leverageChange}
                  ></input>
                )}

                <div
                  className={styles.changeinput}
                  style={{
                    color: color,
                  }}
                  title={"입력 방식을 전환합니다."}
                  onClick={levInputChange}
                >
                  {`x${leverage}`}
                </div>
              </div>
              <div className={styles.submitdiv}>
                <Warning
                  loginWarning={loginWarning}
                  profitWarning={profitWarning}
                  lossWarning={lossWarning}
                  levWarning={levWarning}
                  quanWarning={quanWarning}
                />
                <button
                  className={`${styles.submitbutton} ${
                    isLong ? styles.longsubmit : styles.shortsubmit
                  } ${profitLossErr ? "" : styles.abledbutton}`}
                  disabled={profitLossErr || quanErr || compLoginErr}
                >
                  Submit order
                </button>
              </div>
            </form>
          </div>
          <div
            className={styles.commission}
            title="매 진입시 공제되는 수수료의 비율입니다. 포지션 정리시에는 공제되지 않습니다."
          >
            *Commission : {(commission * 100).toFixed(2)}%
          </div>

          <div className={styles.balancetitle}>BALANCE</div>
          <div className={styles.balanceinterface}>
            <div
              className={styles.balanceroe}
              style={
                balance >= startBalance
                  ? { color: "#26a69a" }
                  : { color: "#ef5350" }
              }
            >
              {Math.floor((10000 * (balance - startBalance)) / startBalance) /
                100}{" "}
              %
            </div>
            <div className={styles.balancebox}>
              <div
                className={styles.balancebody}
                title={`${(balance * 1228.88).toLocaleString("en-US", {
                  maximumFractionDigits: 0,
                })} KRW`}
              >
                <div className={styles.balance}>
                  {balance.toLocaleString("ko-KR", {
                    maximumFractionDigits: 2,
                  })}
                </div>
                <div className={styles.balanceusdt}>USDT</div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
export default OrderInput;
