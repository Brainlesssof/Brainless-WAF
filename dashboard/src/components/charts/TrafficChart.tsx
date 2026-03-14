import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

type TrafficDatum = {
  time: string;
  allowed: number;
  blocked: number;
  throttled: number;
};

type TrafficChartProps = {
  data: TrafficDatum[];
};

export function TrafficChart({ data }: TrafficChartProps) {
  return (
    <div className="traffic-chart">
      <ResponsiveContainer width="100%" height={280}>
        <AreaChart data={data}>
          <defs>
            <linearGradient id="allowedFill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#4ade80" stopOpacity={0.45} />
              <stop offset="95%" stopColor="#4ade80" stopOpacity={0} />
            </linearGradient>
            <linearGradient id="blockedFill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#fb7185" stopOpacity={0.42} />
              <stop offset="95%" stopColor="#fb7185" stopOpacity={0} />
            </linearGradient>
            <linearGradient id="throttledFill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#f59e0b" stopOpacity={0.4} />
              <stop offset="95%" stopColor="#f59e0b" stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid stroke="rgba(255,255,255,0.08)" vertical={false} />
          <XAxis dataKey="time" stroke="rgba(238, 242, 255, 0.58)" tickLine={false} axisLine={false} />
          <YAxis stroke="rgba(238, 242, 255, 0.58)" tickLine={false} axisLine={false} width={36} />
          <Tooltip
            contentStyle={{
              background: "rgba(13, 19, 35, 0.96)",
              border: "1px solid rgba(143, 163, 191, 0.2)",
              borderRadius: "16px",
            }}
          />
          <Area type="monotone" dataKey="allowed" stackId="traffic" stroke="#4ade80" fill="url(#allowedFill)" />
          <Area type="monotone" dataKey="blocked" stackId="traffic" stroke="#fb7185" fill="url(#blockedFill)" />
          <Area type="monotone" dataKey="throttled" stackId="traffic" stroke="#f59e0b" fill="url(#throttledFill)" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}

export default TrafficChart;
