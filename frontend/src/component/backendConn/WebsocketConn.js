// import { useState, useEffect, useRef } from "react";
// import PropTypes from "prop-types";

// function WebsocketConn(mode) {
//   const [tenCharts, setTenCharts] = useState([
//     {
//       onechart: [
//         {
//           close: 0,
//           high: 0,
//           low: 0,
//           open: 0,
//           pairname: "",
//           time: 0,
//           volume: 0,
//         },
//       ],
//       interval: "",
//       backsteps: 0,
//       tradingvalue: 0,
//       name: "",
//     },
//   ]);
//   var result = {};
//   const [socketConnected, setSocketConnected] = useState(false);
//   const [dataLoaded, setDataLoaded] = useState(false);
//   const webSocketUrl = "ws://localhost:80/ws";
//   let ws = useRef(null); // const = 재선언 재할당 불가능, let 재선언 불가능 재할당 가능 => useRef로 초기 할당된 객체에 다른 value를 재할당 하기 위함.

//   useEffect(() => {
//     if (!ws.current) {
//       ws.current = new WebSocket(webSocketUrl);
//       ws.current.onopen = () => {
//         console.log("connected to " + webSocketUrl);
//         setSocketConnected(true);
//       };
//       ws.current.onclose = (error) => {
//         console.log("disconnect from " + webSocketUrl);
//         console.log(error);
//       };
//       ws.current.onerror = (error) => {
//         console.log("connection error " + webSocketUrl);
//         console.log(error);
//       };
//       ws.current.onmessage = (event) => {
//         const data = JSON.parse(event.data);
//         setTenCharts(data);
//         console.log("Websocket Data : ", data);
//         setDataLoaded(true);
//       };
//     }
//     return () => {
//       console.log("Clean up");
//       ws.current.close();
//     };
//   }, []);

//   useEffect(() => {
//     if (socketConnected) {
//       console.log("Mode : ", mode);
//       ws.current.send(JSON.stringify({ message: mode }));
//     }
//   }, [socketConnected]);

//   result.data = tenCharts;
//   result.loaded = dataLoaded;
//   return result;
// }

// export default WebsocketConn;
