﻿//using PhantomInterop.Interfaces;
//using PhantomInterop.Structs.PhantomStructs;
//using PhantomInterop.Structs.MythicStructs;
//using System;
//using System.Collections.Generic;
//using System.IO.Pipes;
//using System.Linq;
//using System.Text;

//namespace PhantomInterop.Serializers
//{
//    public class EncryptedSMBSerializer : JsonHandler, ICryptographySerializer
//    {
//        private ICryptography Cryptor;
        
//        public EncryptedSMBSerializer(ICryptography crypto) : base()
//        {
//            Cryptor = crypto;
//        }

//        public bool UpdateUUID(string uuid)
//        {
//            return Cryptor.UpdateUUID(uuid);
//        }

//        public bool UpdateKey(string key)
//        {
//            return Cryptor.UpdateKey(key);
//        }

//        public string GetUUID()
//        {
//            return Cryptor.GetUUID();
//        }

//        public override string Serialize(object msg)
//        {
//            string jsonMessage = Cryptor.Encrypt(base.Serialize(msg));
//            Type t = msg.GetType();
//            PeerMessage pmsg = new PeerMessage();
//            pmsg.Message = jsonMessage;
//            if (t == typeof(MessageResponse))
//            {
//                pmsg.Type = Enums.PhantomEnums.MessageType.MessageResponse;
//            } else if (t == typeof(CheckinMessage))
//            {
//                pmsg.Type = Enums.PhantomEnums.MessageType.CheckinMessage;
//            } else
//            {
//                throw new Exception($"Invalid message type: {t.Name}");
//            }

//            return base.Serialize(pmsg);
//        }

//        public override T Deserialize<T>(string msg)
//        {
//            PeerMessage pmsg = base.Deserialize<PeerMessage>(msg);
//            string decrypted = Cryptor.Decrypt(pmsg.Message);
//            // do some matching of T to pmsg type and throw exception if not proper.

//        }

//    }
//}
