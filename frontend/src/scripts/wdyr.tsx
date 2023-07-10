/// <reference types="@welldone-software/why-did-you-render" />
import React from 'react';
if (process.env.NEXT_PUBLIC_ENV === 'development') {
  if (typeof window !== 'undefined') {
    const whyDidYouRender = require('@welldone-software/why-did-you-render');
    whyDidYouRender(React, {
      trackAllPureComponents: true,
    });
  }
}
