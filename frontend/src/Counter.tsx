function Counter({ handleIncrement, count }) {
	return (
		<button onClick={handleIncrement} className="counter-button">
			{count}
		</button>
	)
}

export default Counter;
