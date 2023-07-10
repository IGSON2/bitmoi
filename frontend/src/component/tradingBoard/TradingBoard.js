import { useState, useRef, useEffect } from "react";
import styles from "./TradingBoard.module.css";
import ChartRef from "./ChartRef/ChartRef";
import V_Navbar from "../navbar/V_NavBar";
import OrderInput from "./orderInput/OrderInput";
import {
  BsFillArrowRightSquareFill,
  BsFillArrowLeftSquareFill,
} from "react-icons/bs";
import ChartHeader from "./ChartHeader/ChartHeader";
import Loader from "../loader/Loader";

function TradingBoard({ modeHeight, mode }) {
  const [fiveMinutes, setFiveMinutes] = useState();
  const [fifteenMinutes, setFifteenMinutes] = useState();
  const [oneHour, setOneHour] = useState();
  const [fourHour, setFourHour] = useState();
  const [candles, setCandles] = useState([
    {
      pdata: [
        {
          close: 0,
          high: 0,
          low: 0,
          open: 0,
          time: 0,
        },
      ],
      vdata: [
        {
          value: 0,
          time: 0,
          color: "",
        },
      ],
    },
  ]);
  const [resultChart, setResultChart] = useState([
    {
      pdata: [
        {
          close: 0,
          high: 0,
          low: 0,
          open: 0,
          time: 0,
        },
      ],
      vdata: [
        {
          value: 0,
          time: 0,
          color: "",
        },
      ],
    },
  ]);

  const [resultScore, setResultScore] = useState({
    stage: 0,
    name: "",
    leverage: 0,
    entryprice: 0,
    outtime: 0,
    roe: 0,
    pnl: 0,
    commission: 0,
    isliquidated: false,
  });
  const [toolBar, setToolBar] = useState("NonSelected");
  const [loaded, setloaded] = useState(false);
  const [identifier, setIdentifier] = useState("");
  const [headerInterval, setHeaderInterval] = useState("");
  const [index, setIndex] = useState(0);
  const [entryPrice, setEntryPrice] = useState(0);
  const [profitMarker, setProfitMarker] = useState(0);
  const [lossMarker, setLossMarker] = useState(0);
  const [balance, setBalance] = useState(1000);
  const [name, setName] = useState("");
  const [titleaArray, setTitleaArray] = useState([]);
  const [btcRatio, setBtcRatio] = useState(0);
  const [entryTime, setEntryTime] = useState("");
  const [submitOrder, setSubmitOrder] = useState(false);
  const [opened, setOpened] = useState(false);
  const closeButtonDiv = useRef(null);
  const openclosebuttonClick = () => setOpened((current) => !current);
  const [active, setActive] = useState("");
  const getChartData = async (interval) => {
    var jsonData;
    setloaded(false);
    switch (interval) {
      case "init":
        jsonData = await (
          await fetch(
            "http://bitmoi.co.kr:5000/" +
              mode +
              "?names=" +
              titleaArray.join("")
          )
        ).json();
        jsonData.onechart.pdata.reverse();
        jsonData.onechart.vdata.reverse();
        setFiveMinutes();
        setFifteenMinutes();
        setFourHour();
        setOneHour(jsonData.onechart);
        setCandles(jsonData.onechart);
        setIdentifier(jsonData.identifier);
        setName(jsonData.name);
        if (!jsonData.name.includes("STAGE")) {
          setTitleaArray((current) => [...current, jsonData.name + ","]);
        }
        setBtcRatio(jsonData.btcratio);
        setEntryPrice(jsonData.entry_price);
        setEntryTime(jsonData.entrytime);
        setHeaderInterval("1h");
        break;
      case "5m":
        if (fiveMinutes === undefined) {
          jsonData = await (
            await fetch("http://bitmoi.co.kr:5000/interval", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                reqinterval: "5m",
                identifier: identifier,
                mode: mode,
                stage: titleaArray.length,
              }),
            })
          ).json();
          jsonData.onechart.pdata.reverse();
          jsonData.onechart.vdata.reverse();
          setFiveMinutes(jsonData.onechart);
          setCandles(jsonData.onechart);
        } else {
          setCandles(fiveMinutes);
        }
        setHeaderInterval("5m");
        break;
      case "15m":
        if (fifteenMinutes === undefined) {
          jsonData = await (
            await fetch("http://bitmoi.co.kr:5000/interval", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                reqinterval: "15m",
                identifier: identifier,
                mode: mode,
                stage: titleaArray.length,
              }),
            })
          ).json();
          jsonData.onechart.pdata.reverse();
          jsonData.onechart.vdata.reverse();
          setFifteenMinutes(jsonData.onechart);
          setCandles(jsonData.onechart);
        } else {
          setCandles(fifteenMinutes);
        }
        setHeaderInterval("15m");
        break;
      case "1h":
        setCandles(oneHour);
        setHeaderInterval("1h");
        break;
      case "4h":
        if (fourHour === undefined) {
          jsonData = await (
            await fetch("http://bitmoi.co.kr:5000/interval", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                reqinterval: "4h",
                identifier: identifier,
                mode: mode,
                stage: titleaArray.length,
              }),
            })
          ).json();
          jsonData.onechart.pdata.reverse();
          jsonData.onechart.vdata.reverse();
          setFourHour(jsonData.onechart);
          setCandles(jsonData.onechart);
        } else {
          setCandles(fourHour);
        }
        setHeaderInterval("4h");
        break;
    }

    setloaded(true);
  };
  useEffect(() => {
    getChartData("init");
  }, [index]);

  useEffect(() => {
    if (submitOrder) {
      setHeaderInterval("submit");
    }
  }, [submitOrder]);

  window.onkeydown = (e) => {
    setToolBar("NonSelected");
    switch (e.key) {
      case "Shift":
        setActive("ruler");
        break;
      case "Control":
        setActive("mark");
        break;
      case "Alt":
        setActive("horizon");
        break;
    }
  };
  window.onkeyup = () => {
    setActive("");
  };
  return (
    <div className={styles.page}>
      {loaded ? (
        <div className={styles.loadedpage}>
          <div className={styles.top}>
            <ChartHeader
              name={name}
              entryTime={entryTime}
              entryPrice={entryPrice}
              btcRatio={btcRatio}
              getChartData={getChartData}
              headerInterval={headerInterval}
              active={active}
              setActive={setActive}
              setToolBar={setToolBar}
              toolBar={toolBar}
            />
          </div>

          <div className={styles.middle}>
            <div
              className={`${styles.navbar} ${
                opened ? styles.navshow : styles.navclose
              }`}
            >
              {opened ? <V_Navbar /> : null}
            </div>
            <div className={styles.openclosebutton}>
              {opened ? (
                <button
                  onClick={openclosebuttonClick}
                  className={styles.closebutton}
                  title="Close"
                  ref={closeButtonDiv}
                >
                  <BsFillArrowLeftSquareFill />
                </button>
              ) : (
                <button
                  onClick={openclosebuttonClick}
                  className={styles.openbutton}
                  title="Menu"
                >
                  <BsFillArrowRightSquareFill />
                </button>
              )}
            </div>
            <div
              className={`${
                opened ? styles.chartinterface : styles.widerchart
              }`}
            >
              <ChartRef
                candles={candles}
                loaded={loaded}
                submitOrder={submitOrder}
                setSubmitOrder={setSubmitOrder}
                modeHeight={modeHeight}
                opened={opened}
                entryMarker={entryPrice}
                profitMarker={profitMarker}
                lossMarker={lossMarker}
                resultChart={resultChart}
                resultScore={resultScore}
                toolBar={toolBar}
                setToolBar={setToolBar}
                ref={closeButtonDiv}
              />
            </div>
            <div
              className={`${styles.orderInput} ${
                opened ? styles.navshow_orderInput : styles.navclose_orderInput
              }`}
            >
              <OrderInput
                mode={mode}
                name={name}
                index={index}
                opened={opened}
                submitOrder={submitOrder}
                setSubmitOrder={setSubmitOrder}
                entryPrice={entryPrice}
                setProfitMarker={setProfitMarker}
                setLossMarker={setLossMarker}
                identifier={identifier}
                setName={setName}
                setResultChart={setResultChart}
                setCandles={setCandles}
                setIndex={setIndex}
                setResultScore={setResultScore}
                balance={balance}
                setBalance={setBalance}
                setTitleaArray={setTitleaArray}
                entryTime={entryTime}
              />
            </div>
          </div>
        </div>
      ) : (
        <Loader />
      )}
    </div>
  );
}
export default TradingBoard;
