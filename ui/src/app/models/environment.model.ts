export interface EnvironmentModel {
  production: boolean;
  googleMapsApiKey: string;
  beUrl: string;
  firebaseConfig: {
    apiKey: string;
    authDomain: string;
    databaseURL: string;
    projectId: string;
    storageBucket: string;
    messagingSenderId: string;
    appId: string;
    measurementId: string;
  };
}
