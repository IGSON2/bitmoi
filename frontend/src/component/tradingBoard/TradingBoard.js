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
import axiosClient from "../backendConn/axiosClient";
import checkAccessTokenValidity from "../backendConn/checkAccessTokenValidity";

function TradingBoard({ modeHeight, mode, score_id }) {
  const [isLogined, setIsLogined] = useState(false);
  const [userInfo, setUserinfo] = useState();

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
    entry_price: 0,
    end_price: 0,
    out_time: 0,
    roe: 0,
    pnl: 0,
    commission: 0,
    is_liquidated: false,
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
  const [active, setActive] = useState("");

  const openclosebuttonClick = () => setOpened((current) => !current);
  const reqinterval = async (reqinterval, identifier, stage) => {
    const fIdentifier = encodeURIComponent(identifier);
    const reqURL = `/interval?mode=${mode}&reqinterval=${reqinterval}&identifier=${fIdentifier}&stage=${stage}`;
    const response = await axiosClient.get(reqURL);
    return response.data;
  };

  const getChartData = async (interval) => {
    var response;
    setloaded(false);
    switch (interval) {
      case "init":
        var response;
        if (mode === "competition") {
          if (isLogined) {
            response = await axiosClient.get(
              `/competition?names=${titleaArray}`
            );
          } else {
            return;
          }
        } else {
          response = await axiosClient.get(`/practice?names=${titleaArray}`);
        }

        response.data.onechart.pdata.reverse();
        response.data.onechart.vdata.reverse();
        setFiveMinutes();
        setFifteenMinutes();
        setFourHour();
        setOneHour(response.data.onechart);
        setCandles(response.data.onechart);
        setIdentifier(response.data.identifier);
        setName(response.data.name);
        if (!response.data.name.includes("STAGE")) {
          setTitleaArray((current) => [...current, response.data.name + ","]);
        }
        setBtcRatio(response.data.btcratio);
        setEntryPrice(response.data.entry_price);
        setEntryTime(response.data.entrytime);
        setHeaderInterval("1h");
        break;
      case "5m":
        if (fiveMinutes === undefined) {
          const data = reqinterval("5m", identifier, titleaArray.length);
          data.onechart.pdata.reverse();
          data.onechart.vdata.reverse();
          setFiveMinutes(data.onechart);
          setCandles(data.onechart);
        } else {
          setCandles(fiveMinutes);
        }
        setHeaderInterval("5m");
        break;
      case "15m":
        if (fifteenMinutes === undefined) {
          const data = reqinterval("15m", identifier, titleaArray.length);
          data.onechart.pdata.reverse();
          data.onechart.vdata.reverse();
          setFifteenMinutes(data.onechart);
          setCandles(data.onechart);
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
        if (!fourHour) {
          const data = reqinterval("4h", identifier, titleaArray.length);
          data.onechart.pdata.reverse();
          data.onechart.vdata.reverse();
          setFourHour(data.onechart);
          setCandles(data.onechart);
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
  }, [index, isLogined]);

  useEffect(() => {
    if (submitOrder) {
      setHeaderInterval("submit");
    }
  }, [submitOrder]);

  useEffect(() => {
    const verifyToken = async () => {
      const userInfo = await checkAccessTokenValidity();
      if (!userInfo) {
        setIsLogined(false);
        if (mode === "competition") {
          alert("로그인이 필요합니다.");
          window.location.replace("/login");
        }
        setUserinfo({ user_id: "" });
        return;
      }
      setUserinfo(userInfo);
      setIsLogined(true);
    };

    verifyToken();
  }, []);

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
              userInfo={userInfo}
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
                userInfo={userInfo}
                score_id={score_id}
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
