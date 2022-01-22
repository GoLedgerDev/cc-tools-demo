import { Request, Response, NextFunction } from 'express';
import Client from '../createClient';
import invoke from '../../fabric-artifacts/chaincode/invoke';

const invokeTx = (req: Request, res: Response, next: NextFunction) => {
  const client = Client.get();

  let transientMap: Object = {};
  let txRequest: Object = {};
  let collectionsNames: string[];

  // Get transaction name
  const { tx } = req.params;

  // Get collection names
  let collections
  if (req.query['@collections']) {
    collections = req.query['@collections']
    try {
      const collectionsJSON = Buffer.from(collections, 'base64').toString('utf-8');
      collectionsNames = JSON.parse(collectionsJSON);
      if (!Array.isArray(collectionsNames)) {
        throw '';
      }
    } catch (err) {
      return res.status(400).send('the @collections query parameter must be a base64-encoded JSON array of strings');
    }
  } else if (req.query['collections']) {
    if (Array.isArray(req.query['collections'])) {
      collectionsNames = req.query['collections'] as string[]
    } else {
      collectionsNames = [req.query['collections']] as string[]
    }
  }

  // Handle req.body
  for (const key in req.body) {
    if (key.charAt(0) == '~') {
      transientMap[key.slice(1, key.length)] = req.body[key];
    } else {
      txRequest[key] = req.body[key];
    }
    
  }

  const txArgs = [JSON.stringify(txRequest)];
  const transientRequest = new Buffer(JSON.stringify(transientMap));

  invoke(client, tx, txArgs, transientRequest, collectionsNames)
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

export default invokeTx;
