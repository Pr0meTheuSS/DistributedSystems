import { useEffect, useRef, useState } from 'react';
import axios from 'axios';
import './App.css';

type Request = {
  id: string,
  hash: string,
  progress: number,
  status: string,
};

function App() {
  const [requests, setRequests] = useState<Request[]>([]);
  const [hash, setHash] = useState("");
  const [maxLength, setMaxLength] = useState("");
  const [alphabet, setAlphabet] = useState("");
  const [subscribeId, setSubscribeId] = useState("");
  const requestIdsRef = useRef(new Set());

  const presets = [
    { label: "Lowercase a-z", value: "abcdefghijklmnopqrstuvwxyz" },
    { label: "Uppercase A-Z", value: "ABCDEFGHIJKLMNOPQRSTUVWXYZ" },
    { label: "Numbers 0-9", value: "0123456789" },
    { label: "Hexadecimal", value: "0123456789abcdef" },
    { label: "AlphaNumeric", value: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" },
    { label: "Symbols", value: "!@#$%^&*()_+-=[]{}|;:',.<>/?`~" },
  ];

  useEffect(() => {
    const interval = setInterval(() => {
      requests.forEach(async (req: Request) => {
        try {
          const res = await axios.get(`http://localhost:9091/api/hash-cracks/${req.id}/progress`);
          setRequests((prevRequests) =>
            prevRequests.map((r: Request) =>
              r.id === req.id
                ? {
                    ...r,
                    progress: res.data.progress_in_percents,
                    status: res.data.progress_in_percents === 100 ? "DONE" : "IN_PROGRESS",
                    hash: res.data.hash,
                  }
                : r
            )
          );
        } catch (err) {
          console.error(`Error for fetching progress for request ID: ${req.id}:`, err);
        }
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [requests]);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();

    if (!hash || !maxLength || !alphabet) {
      alert("Need to fill all inputs");
      return;
    }

    try {
      const response = await axios.post('http://localhost:9091/api/hash-cracks', {
        hash,
        length: Number(maxLength),
        alphabet,
      });

      const newId = response.data.id;

      if (!requestIdsRef.current.has(newId)) {
        requestIdsRef.current.add(newId);
        setRequests((prev) => [
          ...prev,
          {
            id: newId,
            hash: "",
            progress: 0,
            status: "IN_PROGRESS",
          },
        ]);
      }
    } catch (error) {
      console.error("Request erro:", error);
      alert("Request error");
    }
  };

  const handleSubscribe = async () => {
    if (!subscribeId) {
      alert("Enter ID for subscribe");
      return;
    }

    try {
      const res = await axios.get(`http://localhost:9091/api/hash-cracks/${subscribeId}/progress`);

      if (!requestIdsRef.current.has(subscribeId)) {
        requestIdsRef.current.add(subscribeId);
        setRequests((prev) => [
          ...prev,
          {
            id: subscribeId,
            hash: res.data.hash,
            progress: res.data.progress_in_percents,
            status: res.data.progress_in_percents === 100 ? "DONE" : "IN_PROGRESS",
          },
        ]);
      }

      setSubscribeId("");
    } catch (err) {
      console.error("Error while subscribing:", err);
      alert("Error while subscribing. Check request ID is correct.");
    }
  };

  return (
    <div className='input-form-wrapper'>
      <form className='input-form' onSubmit={handleSubmit}>
        <label>Hash</label>
        <input
          className='hash-input'
          type='text'
          placeholder='hash'
          value={hash}
          onChange={(e) => setHash(e.target.value)}
        />

        <label>Max Length</label>
        <input
          className='max-length-input'
          type='number'
          placeholder='max length'
          value={maxLength}
          onChange={(e) => setMaxLength(e.target.value)}
        />

        <label>Preset Alphabet</label>
        <select
          className='alphabet-select'
          onChange={(e) => setAlphabet(e.target.value)}
          value=""
        >
          <option value="">-- Select preset --</option>
          {presets.map((preset) => (
            <option key={preset.label} value={preset.value}>
              {preset.label}
            </option>
          ))}
        </select>

        <label>Custom Alphabet</label>
        <input
          className='alphabet-input'
          type='text'
          value={alphabet}
          onChange={(e) => setAlphabet(e.target.value)}
          placeholder='Enter custom alphabet'
        />

        <button type='submit'>Send Request</button>
      </form>

      <div className='subscribe-section'>
        <label>Подписаться на request ID</label>
        <input
          type='text'
          placeholder='Enter request ID'
          value={subscribeId}
          onChange={(e) => setSubscribeId(e.target.value)}
        />
        <button onClick={handleSubscribe}>Subscribe</button>
      </div>

      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Hash</th>
            <th>Progress</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {requests.map((request) => (
            <tr key={request.id}>
              <td>{request.id}</td>
              <td>{request.hash}</td>
              <td>{request.progress?.toFixed(2)}%</td>
              <td>{request.status}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;
