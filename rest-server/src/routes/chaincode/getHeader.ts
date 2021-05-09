import { Request, Response, NextFunction } from 'express';
import Client from '../createClient';
import query from '../../fabric-artifacts/chaincode/query';

const getHeader = (req: Request, res: Response, next: NextFunction) => {
  const client = Client.get();

  query(client, 'getHeader', [])
    .then(response => {
      return res.send(response);
    })
    .catch(err => {
      console.error(err);
      if (err && err.status) {
        return res.status(err.status).send(err.message);
      }
      return res.status(500).send(err);
    });
};

export default getHeader;
