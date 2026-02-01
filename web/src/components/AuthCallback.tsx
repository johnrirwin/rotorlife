import { useEffect, useState } from 'react';
import { storeTokens, getCurrentUser } from '../authApi';

export function AuthCallback() {
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(true);

  useEffect(() => {
    const handleCallback = async () => {
      // Get tokens from URL fragment
      const hash = window.location.hash.substring(1);
      const params = new URLSearchParams(hash);
      
      const accessToken = params.get('access_token');
      const refreshToken = params.get('refresh_token');

      // Check for error in query params
      const searchParams = new URLSearchParams(window.location.search);
      const errorParam = searchParams.get('error');
      const errorDesc = searchParams.get('error_description');

      if (errorParam) {
        setError(errorDesc || errorParam);
        setIsProcessing(false);
        return;
      }

      if (!accessToken || !refreshToken) {
        setError('Missing authentication tokens');
        setIsProcessing(false);
        return;
      }

      try {
        // Store tokens using the authApi function
        storeTokens({
          accessToken,
          refreshToken,
          tokenType: 'Bearer',
          expiresIn: 3600,
        });
        
        // Verify tokens work by fetching user - this ensures tokens are valid
        await getCurrentUser();
        
        // Small delay to ensure localStorage is synced, then redirect
        await new Promise(resolve => setTimeout(resolve, 100));
        window.location.replace('/');
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Authentication failed');
        setIsProcessing(false);
      }
    };

    handleCallback();
  }, []);

  if (error) {
    return (
      <div className="fixed inset-0 bg-gray-900 flex items-center justify-center">
        <div className="bg-gray-800 rounded-lg p-6 max-w-md mx-4">
          <h2 className="text-xl font-bold text-white mb-4">Authentication Error</h2>
          <p className="text-red-400 mb-4">{error}</p>
          <button
            onClick={() => window.location.href = '/'}
            className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            Return to Home
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-gray-900 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
        <p className="text-white">{isProcessing ? 'Completing sign in...' : 'Redirecting...'}</p>
      </div>
    </div>
  );
}
