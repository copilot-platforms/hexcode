import { BarChart } from '@mui/x-charts';
import preactLogo from '../../assets/preact.svg';
import './style.css';

export function Home() {
	return (
		<div class="home">
			<a href="https://preactjs.com" target="_blank">
				<img src={preactLogo} alt="Preact logo" height="160" width="160" />
			</a>
			<BarChart  xAxis={[
    {
      id: 'barCategories',
      data: ['bar A', 'bar B', 'bar C'],
      scaleType: 'band',
    },
	
  ]}
  series={[
    {
      data: [2, 5, 3],
    },
  ]}
  width={500}
  height={300}></BarChart>
		
		</div>
	);
}

function Resource(props) {
	return (
		<a href={props.href} target="_blank" class="resource">
			<h2>{props.title}</h2>
			<p>{props.description}</p>
		</a>
	);
}
