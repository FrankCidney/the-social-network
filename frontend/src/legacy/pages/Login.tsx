import * as React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';

const Login = () => {
  const [isLoading, setIsLoading] = React.useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    // Simulate login
    setTimeout(() => setIsLoading(false), 1500);
  };

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h2 className="text-2xl font-bold tracking-tight text-text-main">Welcome back</h2>
        <p className="text-sm text-gray-500">
          Enter your credentials to access your account
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input 
          label="Email or Nickname" 
          placeholder="name@example.com" 
          type="text"
          required
        />
        <div className="space-y-1">
          <Input 
            label="Password" 
            placeholder="••••••••" 
            type="password"
            required
          />
          <div className="flex justify-end">
            <button type="button" className="text-xs text-primary hover:underline">
              Forgot password?
            </button>
          </div>
        </div>
        
        <Button 
          type="submit" 
          className="w-full" 
          disabled={isLoading}
        >
          {isLoading ? 'Signing in...' : 'Sign In'}
        </Button>
      </form>

      <div className="text-center text-sm">
        <span className="text-gray-500">Don't have an account? </span>
        <Link to="/register" className="font-medium text-primary hover:underline">
          Create an account
        </Link>
      </div>
    </div>
  );
};

export default Login;
