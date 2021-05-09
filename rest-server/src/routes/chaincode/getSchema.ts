import { Request, Response, NextFunction } from 'express';
import Client from '../createClient';
import query from '../../fabric-artifacts/chaincode/query';

const getSchema = (req: Request, res: Response, next: NextFunction) => {
  const client = Client.get();

  const { asset } = req.params;

  const args = asset ? [JSON.stringify({ assetType: asset })] : [];

  query(client, 'getSchema', args)
    .then((response) => {
      return res.send(response);
    })
    .catch((err) => {
      console.error(err);
      if (err && err.status) {
        return res.status(err.status).send(err.message);
      }
      return res.status(500).send(err);
    });
};

export default getSchema;
