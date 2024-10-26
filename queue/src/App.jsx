import { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [progressData, setProgressData] = useState({}); // Store progress data for rendering
  const [completedTasks, setCompletedTasks] = useState({}); // Store completed tasks
  const [colorMapping, setColorMapping] = useState({}); // Store mapping of UUID to color
  const [loading, setLoading] = useState(false); // Loading state

  // Fetch progress data from the server every second
  useEffect(() => {
    setLoading(true);
    const intervalId = setInterval(async () => {
      try {
        const response = await fetch('http://localhost:3333/taskmon');
        const data = await response.json();

        const newCompletedTasks = {};
        for (const taskID in data) {
          if (data[taskID].state === "Complete") {
            newCompletedTasks[taskID] = data[taskID];
          }
        }
        setCompletedTasks((prev) => ({ ...prev, ...newCompletedTasks })); // Store completed tasks

        // Remove completed tasks from progressData
        const updatedProgressData = { ...data };
        for (const taskID in newCompletedTasks) {
          delete updatedProgressData[taskID];
        }

        setProgressData(updatedProgressData);

      } catch (error) {
        console.error('Error fetching tasks:', error);
      } finally {
        setLoading(false);
      }
    }, 100); 

    return () => clearInterval(intervalId);
  }, []);

  const spawnTask = async () => {
    const portMap = {
      0: 8080,
      1: 8081,
      2: 8082,
    };

    const colors = {
      0: 'yellow',
      1: 'blue',
      2: 'pink',
    };

    // Generate a random number between 0 and 2
    const num = Math.floor(Math.random() * 3);
    const port = portMap[num];
    const color = colors[num];

    try {
      const response = await fetch(`http://localhost:${port}/work`, {
        method: 'POST', // Send a POST request
      });
      const result = await response.json();
      const taskID = result.taskID; // Extract taskID from response
      console.log('Spawned Task ID:', taskID);

      setColorMapping((prev) => ({
        ...prev,
        [taskID]: color,
      }));

      setProgressData((prev) => ({
        ...prev,
        [taskID]: {
          ...result,
        },
      }));



    } catch (error) {
      console.error('Error spawning task:', error);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold underline mb-4">Task Tracker</h1>

      <div className="flex justify-between">
        <div className="w-1/2 pr-2">
          <h2 className="text-lg font-bold mt-4">Running Tasks</h2>
          {Object.keys(progressData).map((taskID) => (
            <TaskCard key={taskID} taskID={taskID} progressData={progressData} color={colorMapping[taskID]} />
          ))}
        </div>

        <div className="w-1/2 pl-2">
          <h2 className="text-lg font-bold mt-4">Completed Tasks</h2>
          {Object.keys(completedTasks).map((taskID) => (
            <TaskCard key={taskID} taskID={taskID} progressData={completedTasks} />
          ))}
        </div>
      </div>

      <div className="mt-4">
        <button
          className="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
          onClick={spawnTask} 
        >
          Hit (Spawn Random Task)
        </button>
      </div>
    </div>
  );
}

const TaskCard = ({ taskID, progressData, color }) => {
  const task = progressData[taskID]; 

  console.log(`${color || "lime"} `)
  return (


    <div className="mt-4 p-4 border rounded shadow-md">
      <h2 className="font-semibold">{task?.name || 'Unknown Task'}</h2>
      <p className="text-sm text-gray-600">ID: {taskID}</p>
      <p className="text-sm text-gray-600">State: {task?.state || 'Unknown'}</p>
      <div className="mt-2">
        <div className="bg-gray-200 rounded-full h-2">
          <div
            className="h-2 rounded-full"
            style={{
              width: `${task?.percentageComplete * 100 || 0}%`,
              backgroundColor: `${color || "lime"} `,
            }}
          />
        </div>
        <p className="text-sm text-gray-600">
          Progress: {((task?.percentageComplete || 0) * 100).toFixed(2)}%
        </p>
      </div>
    </div>
  );
};

export default App;
