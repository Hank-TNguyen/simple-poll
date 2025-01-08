import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

function ViewPoll() {
  const { pollId } = useParams();
  const [poll, setPoll] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchPoll() {
      try {
        setLoading(true);
        setError(null);

        // Replace with your actual backend endpoint (port, route, etc.)
        const response = await fetch(`http://localhost:3000/api/polls/${pollId}`);
        if (!response.ok) {
          throw new Error(`Failed to fetch poll with ID ${pollId}`);
        }

        const data = await response.json();
        setPoll(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchPoll();
  }, [pollId]);

  if (loading) return <div>Loading poll...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!poll) return <div>Poll not found.</div>;

  return (
    <div>
      <h1>{poll.title}</h1>
      <p>{poll.description}</p>

      {poll.questions && poll.questions.length > 0 ? (
        <div>
          {poll.questions.map((question) => (
            <div key={question.id} style={{ margin: '1rem 0' }}>
              <h3>{question.text}</h3>
              <ul>
                {question.choices?.map((choice) => (
                  <li key={choice.id}>
                    {/* Use choice.choice_text since that's what the backend returns */}
                    {choice.choice_text || '(No text provided)'}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      ) : (
        <p>No questions available for this poll.</p>
      )}
    </div>
  );
}

export default ViewPoll;