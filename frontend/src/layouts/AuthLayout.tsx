import * as React from 'react';
import { Outlet } from 'react-router-dom';
import { motion } from 'framer-motion';

const AuthLayout = () => {
  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-canvas font-sans">
      {/* Visual Narrative Side */}
      <div className="hidden md:flex md:w-1/2 bg-primary p-12 flex-col justify-between text-white relative overflow-hidden">
        <div className="z-10">
          <h1 className="text-4xl font-bold tracking-tight">Social Network</h1>
          <p className="mt-4 text-indigo-100 max-w-md text-lg">
            Connect with your world, share your moments, and maintain your privacy—all in one place.
          </p>
        </div>
        
        {/* Abstract Bento Shape Decoration */}
        <div className="absolute top-0 right-0 -translate-y-1/2 translate-x-1/2 w-[600px] h-[600px] bg-indigo-500 rounded-full opacity-20 blur-3xl" />
        <div className="absolute bottom-0 left-0 translate-y-1/4 -translate-x-1/4 w-[400px] h-[400px] bg-white rounded-full opacity-10 blur-2xl" />

        <div className="z-10">
          <div className="flex -space-x-2 overflow-hidden">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="inline-block h-10 w-10 rounded-full ring-2 ring-primary bg-indigo-400" />
            ))}
          </div>
          <p className="mt-4 text-sm text-indigo-100">
            Join thousands of users sharing their stories today.
          </p>
        </div>
      </div>

      {/* Form Side */}
      <div className="flex-1 flex items-center justify-center p-6 md:p-12">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="w-full max-w-[400px] space-y-8"
        >
          <Outlet />
        </motion.div>
      </div>
    </div>
  );
};

export default AuthLayout;
