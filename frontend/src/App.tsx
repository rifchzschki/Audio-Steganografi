import { useEffect, useState } from 'react';
import './App.css';
import type { BaseResponseMessage } from './models/response';

function App() {
  const [data, setData] = useState<BaseResponseMessage | null>(null);

  useEffect(() => {
    fetchData();
  }, []);
  const fetchData = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/hello');
      const data = await response.json();
      setData(data);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };
  return (
    <>
      <h1 className="text-xl font-bold text-red-500">Audio Steganography</h1>
      <p className="text-blue-500">{data?.message}</p>
    </>
  );
}

export default App;
