import React from 'react';

const useMount = (callback: () => void) => {
  const mounted = React.useRef(false);

  React.useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;
      callback();
    }
  }, [callback]);
};

export default useMount;
