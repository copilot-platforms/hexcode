import { BarChart, LineChart, PieChart } from "@mui/x-charts";
import { useComputed, useSignal } from "@preact/signals";
import { useEffect } from "preact/hooks";
import eventData from "./data.json";
import "./style.css";

const valueFormatter = (date: Date) =>
  date.getHours() === 0
    ? date.toLocaleDateString("fr-FR", {
        month: "2-digit",
        day: "2-digit",
      })
    : date.toLocaleTimeString("fr-FR", {
        hour: "2-digit",
      });

export function Home() {
  const makeGraphs = (events) =>
    events.data.map((event) => {
      switch (event.type) {
        case "pie":
          return (
            <div>
              <h3>{event.title}</h3>
              <PieChart
                series={[
                  {
                    data: event.data.map((chunk, index) => {
                      return {
                        id: index,
                        value: chunk.count,
                        label: chunk.label,
                      };
                    }),
                  },
                ]}
                width={500}
                height={200}
                className="pie-container"
              />
            </div>
          );

        case "line":
          return (
            <div>
              <h3>{event.title}</h3>
              <LineChart
                xAxis={[
                  {
                    data: event.data.map((chunk, i) => new Date(2023, 10, i)),
                    scaleType: "time",
                    valueFormatter: valueFormatter,
                  },
                ]}
                series={[
                  {
                    data: event.data.map((chunk) => chunk.count),
                  },
                ]}
                width={500}
                height={300}
              />
            </div>
          );
        case "bar-single":
          return (
            <div>
              <h3>{event.title}</h3>
              <BarChart
                xAxis={[
                  {
                    scaleType: "band",
                    data: event.data.map((chunk, i) => chunk.label),
                  },
                ]}
                series={[
                  {
                    data: event.data.map((chunk) => chunk.count),
                  },
                ]}
                width={1200}
                height={400}
              />
            </div>
          );
        case "bar-multiple":
          const series = [...new Set(event.data.map((e) => e.label))];

          const seriesData = series.map((s) => {
            return {
              stack: "total",
              data: event.data.map((e) => (e.label === s ? e.count : 0)),
              label: s,
            };
          });

          const data = [...new Set(event.data.map((chunk, i) => chunk.key))];
          console.log(seriesData);
          console.log(data);

          return (
            <div>
              <h3>{event.title}</h3>
              <BarChart
                xAxis={[
                  {
                    scaleType: "band",
                    data: data,
                  },
                ]}
                series={seriesData}
                width={1200}
                height={400}
              />
            </div>
          );
      }
    });

  const events = useSignal(eventData || { data: [] });
  const graphs = useComputed(() => makeGraphs(events.value));

  async function loadData() {
    const response = await fetch("/data");
    const data = await response.json();
    events.value = data;
  }

  useEffect(() => {
    loadData().then(() => console.log("done"));
  }, []);

  return <div class="home graph-container">{graphs.value}</div>;
}

function Resource(props) {
  return (
    <a href={props.href} target="_blank" class="resource">
      <h2>{props.title}</h2>
      <p>{props.description}</p>
    </a>
  );
}
