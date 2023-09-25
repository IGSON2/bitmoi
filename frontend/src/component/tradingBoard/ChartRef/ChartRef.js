import React, { useEffect, useRef, useState, forwardRef } from "react";
import { createChart, CrosshairMode, LineStyle } from "lightweight-charts";
import Timeformatter from "../../Timeformatter/Timeformatter";

const ChartRef = forwardRef((props, ref) => {
  const backgroundcolor = "#555860";
  const linecolor = "#ffffff";
  const chartContainerRef = useRef();
  const candleSeriesRef = useRef();
  const volumeSeriesRef = useRef();
  const chartRef = useRef();
  var ruler = {};
  var rulerFrom = [];
  var rulerTo = [];
  var marks = [];
  var pricelines = [];
  const [visibleRange, setVisibleRange] = useState();
  useEffect(() => {
    const handleResize = () => {
      chart.applyOptions({
        width: chartContainerRef.current.clientWidth,
        height: document.body.clientHeight * props.modeHeight,
      });
    };
    const chart = createChart(chartContainerRef.current, {
      width: chartContainerRef.current.clientWidth,
      height: document.body.clientHeight * props.modeHeight,
      crosshair: {
        mode: CrosshairMode.Normal,
        vertLine: {
          color: "#D6DCDE",
        },
        horzLine: {
          color: "#D6DCDE",
        },
      },
      layout: {
        textColor: backgroundcolor,
        backgroundColor: linecolor,
      },
      grid: {
        vertLines: {
          color: "#D6DCDE",
          style: LineStyle.Dotted,
        },
        horzLines: {
          color: "#D6DCDE",
          style: LineStyle.Dotted,
        },
      },
      timeScale: {
        timeVisible: true,
        secondsVisible: false,
      },
      localization: {
        dateFormat: "yyyy-MM-dd",
      },
    });
    chartRef.current = chart;

    const candleSeries = chart.addCandlestickSeries();
    const volumeSeries = chart.addHistogramSeries({
      priceFormat: {
        type: "volume",
      },
      priceScaleId: "",
      scaleMargins: {
        top: 0.8,
        bottom: 0,
      },
    });
    candleSeries.applyOptions({
      priceFormat: {
        precision: 4,
        minMove: 0.0001,
      },
    });
    if (props.loaded) {
      candleSeries.setData(props.candles.pdata);
      volumeSeries.setData(props.candles.vdata);
      candleSeriesRef.current = candleSeries;
      volumeSeriesRef.current = volumeSeries;
      if (props.toolBar === "NonSelected") {
        var range = 200;
        const length = props.candles.pdata.length;
        if (length < 200) {
          range = length;
        }
        chart.timeScale().setVisibleRange({
          from: props.candles.pdata[length - range].time,
          to: props.candles.pdata[length - 1].time,
        });
      }
    }

    chart.timeScale().subscribeVisibleTimeRangeChange((param) => {
      setVisibleRange(param);
    });
    if (props.toolBar !== "NonSelected") {
      chart.timeScale().setVisibleRange(visibleRange);
    }

    const clickHandler = (param) => {
      if (window.event.shiftKey || props.toolBar === "ruler") {
        if (!ruler.first) {
          rulerFrom.map((from) => {
            candleSeries.removePriceLine(from);
          });
          rulerTo.map((to) => {
            candleSeries.removePriceLine(to);
          });
          rulerFrom = [];
          rulerTo = [];
          ruler.first = candleSeries.coordinateToPrice(param.point.y);
          ruler.second = undefined;
          rulerFrom = [
            ...rulerFrom,
            candleSeries.createPriceLine({
              price: ruler.first,
              color: "rgb(51, 61, 121)",
              lineWidth: 2,
              lineStyle: LineStyle.Solid,
              axisLabelVisible: true,
              title: "From",
            }),
          ];
        } else {
          if (!ruler.second) {
            ruler.second = candleSeries.coordinateToPrice(param.point.y);
            rulerTo = [
              ...rulerTo,
              candleSeries.createPriceLine({
                price: ruler.second,
                color: "rgb(51, 61, 121)",
                lineWidth: 2,
                lineStyle: LineStyle.Solid,
                axisLabelVisible: true,
                title:
                  ((100 * (ruler.second - ruler.first)) / ruler.first).toFixed(
                    2
                  ) + "%",
              }),
            ];
            ruler = {};
          }
        }
      } else if (window.event.ctrlKey || props.toolBar === "mark") {
        const dateformat = Timeformatter(
          chart.timeScale().coordinateToTime(param.point.x) * 1000,
          true
        );
        marks = [
          ...marks,
          {
            time: chart.timeScale().coordinateToTime(param.point.x),
            position: "aboveBar",
            color: "orange",
            shape: "arrowDown",
            id: "id4",
            text: dateformat,
            size: 2,
          },
        ];
        candleSeries.setMarkers(marks);
      } else if (window.event.altKey || props.toolBar === "horizon") {
        pricelines = [
          ...pricelines,
          candleSeries.createPriceLine({
            price: candleSeries.coordinateToPrice(param.point.y),
            color: "#4c525e",
            lineWidth: 2,
            lineStyle: LineStyle.Solid,
            axisLabelVisible: true,
          }),
        ];
      } else {
        rulerFrom.map((from) => {
          candleSeries.removePriceLine(from);
        });
        rulerTo.map((to) => {
          candleSeries.removePriceLine(to);
        });
        rulerFrom = [];
        rulerTo = [];
        marks = [];
        ruler = {};
        candleSeries.setMarkers([]);
        pricelines.map((line) => {
          candleSeries.removePriceLine(line);
        });
        props.setToolBar("NonSelected");
      }
    };

    if (props.toolBar === "escape") {
      rulerFrom.map((from) => {
        candleSeries.removePriceLine(from);
      });
      rulerTo.map((to) => {
        candleSeries.removePriceLine(to);
      });
      rulerFrom = [];
      rulerTo = [];
      marks = [];
      candleSeries.setMarkers([]);
      pricelines.map((line) => {
        candleSeries.removePriceLine(line);
      });
      props.setToolBar("ruler");
    }

    chart.subscribeClick(clickHandler);

    const resizeByClick = () => {
      if (props.opened) {
        setTimeout(handleResize, 1000);
      }
    };

    window.addEventListener("resize", handleResize);
    if (ref.current != null) {
      ref.current.addEventListener("click", resizeByClick);
    }

    return () => {
      window.removeEventListener("resize", handleResize);
      chart.remove();
    };
  }, [props.candles, props.toolBar]);

  useEffect(() => {
    if (candleSeriesRef.current === null) {
      return;
    }

    if (props.submitOrder && props.loaded) {
      var i = 0;
      candleSeriesRef.current.createPriceLine({
        price: props.resultScore.entry_price,
        color: "rgb(51, 61, 121)",
        lineWidth: 2,
        lineStyle: LineStyle.Dotted,
        axisLabelVisible: true,
        title: "Entry price",
      });
      candleSeriesRef.current.createPriceLine({
        price: props.resultScore.profit_price,
        color: "rgb(53, 182, 169)",
        lineWidth: 2,
        lineStyle: LineStyle.Dotted,
        axisLabelVisible: true,
        title: "Take profit",
      });
      candleSeriesRef.current.createPriceLine({
        price: props.resultScore.loss_price,
        color: "rgb(238, 103, 101)",
        lineWidth: 2,
        lineStyle: LineStyle.Dotted,
        axisLabelVisible: true,
        title: "Stop loss",
      });
      candleSeriesRef.current.setMarkers([
        {
          time: props.candles.pdata[props.candles.pdata.length - 1].time,
          position: "aboveBar",
          color: "rgb(51, 61, 121)",
          shape: "arrowDown",
          size: 1.5,
        },
      ]);
      const interval = setInterval(() => {
        const nextP = {
          time: props.resultChart.pdata[i].time,
          open: props.resultChart.pdata[i].open,
          high: props.resultChart.pdata[i].high,
          low: props.resultChart.pdata[i].low,
          close: props.resultChart.pdata[i].close,
        };
        const nextV = {
          time: props.resultChart.vdata[i].time,
          value: props.resultChart.vdata[i].value,
          color: props.resultChart.vdata[i].color,
        };

        chartRef.current.timeScale().setVisibleRange({
          from: props.candles.pdata[props.candles.pdata.length - 100].time,
          to: props.resultChart.pdata[i].time,
        });

        chartRef.current.applyOptions({
          timeScale: {
            rightOffset: 5,
          },
        });

        candleSeriesRef.current.update(nextP);
        volumeSeriesRef.current.update(nextV);

        if (props.resultScore.out_time - 1 == i) {
          if (props.resultScore.pnl > 0) {
            candleSeriesRef.current.setMarkers([
              {
                time: props.candles.pdata[props.candles.pdata.length - 1].time,
                position: "aboveBar",
                color: "rgb(51, 61, 121)",
                shape: "arrowDown",
                size: 1.5,
              },
              {
                time: props.resultChart.pdata[i].time,
                position: "aboveBar",
                color: "rgb(38,166,154)",
                shape: "arrowDown",
                size: 1.5,
              },
            ]);
          } else {
            candleSeriesRef.current.setMarkers([
              {
                time: props.candles.pdata[props.candles.pdata.length - 1].time,
                position: "aboveBar",
                color: "rgb(51, 61, 121)",
                shape: "arrowDown",
                size: 1.5,
              },
              {
                time: props.resultChart.pdata[i].time,
                position: "aboveBar",
                color: "rgb(239,83,80)",
                shape: "arrowDown",
                size: 1.5,
              },
            ]);
          }
          clearInterval(interval);
          props.setSubmitOrder(false);
        }

        i++;
      }, 100);
      setTimeout(function () {
        clearInterval(interval);
        props.setSubmitOrder(false);
      }, 2400);
    } else {
      return;
    }
  }, [props.resultScore]);

  return <div ref={chartContainerRef} />;
});

export default ChartRef;
