
from argparse import ArgumentParser
import itertools
import numpy as np
import random
import tensorflow as tf
import time
from tensorflow.keras import Input, Model
from tensorflow.keras.layers import Concatenate, Lambda, Layer
from tensorflow.keras.layers.experimental.preprocessing import StringLookup, TextVectorization

class LookupLayer(Layer):
    def __init__(self, keys, values, default_value, name=None, table_name=None):
        self.keys = keys 
        self.values = values
        self.default_value = default_value
        self.table_name = table_name
        super().__init__(name=name, trainable=False)

    def get_config(self):
        return dict()

    def build(self, input_shape):
        keys_tensor = tf.constant(self.keys)
        values_tensor = tf.constant(self.values)
        kv_init = tf.lookup.KeyValueTensorInitializer(keys=keys_tensor, values=values_tensor)
        table = tf.lookup.StaticHashTable(kv_init, default_value=self.default_value, name=self.table_name)
        self.table = table

    def call(self, inputs):
        return self.table.lookup(inputs)

if __name__ == '__main__':
    ap = ArgumentParser(description='Generate TensorFlow 2.4.x models.')

    # incomplete, docs provide a lot more entropy
    ap.add_argument('--seed', type=int)

    ap.add_argument('--vocab-size', type=int, default=3)
    ap.add_argument('--max-test-vector-length', type=int, default=8)
    ap.add_argument('--num-test-samples', type=int, default=4)

    ap.add_argument('--no-vectorization', action='store_true')
    ap.add_argument('--no-string-lookup', action='store_true')
    ap.add_argument('--no-float-input', action='store_true')

    ap.add_argument('--no-int-out', action='store_true')
    ap.add_argument('--no-float-out', action='store_true')

    ap.add_argument('--no-lookup-layer', action='store_true', help='for broken_case_v0_5_0 model')
    ap.add_argument('--no-keyed-out', action='store_true', help='for testing v0.5.0')

    ap.add_argument('--no-slow', action='store_true', help='for testing thread limit')
    ap.add_argument('--slow-repeat', type=int, default=10, help='looping math operations for slow model')
    ap.add_argument('--slow-depth', type=int, default=4096)
    ap.add_argument('--slow-depth-exp', type=int, default=2)

    pa, _ = ap.parse_known_args()

    if pa.seed:
        random.seed(pa.seed)

    vocab_size = pa.vocab_size
    vocab_data = [chr(ord('a') + i) for i in range(vocab_size)]
    print(vocab_data)