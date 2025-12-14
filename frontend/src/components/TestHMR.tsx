
import React from 'react';

/**
 * Test component for Hot Module Replacement verification
 */
export const TestHMR: React.FC = () => {
  return (
    <div style={{ 
      padding: '20px', 
      backgroundColor: '#f0f0f0', 
      border: '2px solid #007acc',
      borderRadius: '8px',
      margin: '10px'
    }}>
      <h3>HMR Test Component - Version 1</h3>
      <p>This component tests hot module replacement functionality.</p>
      <p>Current time: {new Date().toLocaleTimeString()}</p>
    </div>
  );
};

export default TestHMR;

// Enable HMR for this component
if (import.meta.hot) {
  import.meta.hot.accept();
}
